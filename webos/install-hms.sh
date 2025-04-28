#!/bin/sh

PREFIX=/usr
LIBDIR=$PREFIX/lib

. "$LIBDIR/libalpine.sh"
. "$LIBDIR/dasd-functions.sh"

ROOTFS=${ROOTFS:-ext4}
BOOTFS=${BOOTFS:-ext4}
VARFS=${VARFS:-ext4}
DISKLABEL=${DISKLABEL:-dos}

# default location for mounted root
SYSROOT=${SYSROOT:-/mnt}

# machine arch
ARCH=$(apk --print-arch)
BOOTLOADER=${BOOTLOADER:-syslinux}

if [ -t 1 ]; then
	COLCYAN="\e[36m"
	COLWHITE="\e[97m"
	COLRESET="\e[0m"
else
	COLCYAN=""
	COLWHITE=""
	COLRESET=""
fi

print_heading1() {
	printf "${COLCYAN}%s${COLRESET}\n" "$1"
}

print_heading2() {
	printf "${COLWHITE}%s${COLRESET}\n" "$1"
}

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat >"$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

partition_id() {
	local id

	if [ "$DISKLABEL" != "dos" ]; then
		die "Partition label \"$DISKLABEL\" is not supported!"
	fi

	case "$1" in
	swap) id=82 ;;
	linux) id=83 ;;
	prep) id=41 ;;
	*) die "Partition id \"$1\" is not supported!" ;;
	esac

	echo $id
}

init_chroot_mounts() {
	local mnt="$1" i=
	for i in proc dev; do
		mkdir -pv "$mnt"/$i
		$MOCK mount --bind /$i "$mnt"/$i
	done
}

cleanup_chroot_mounts() {
	local mnt="$1" i=
	for i in proc dev; do
		umount "$mnt"/$i
	done
}

find_mount_dev() {
	local mnt="$1"
	awk "\$2 == \"$mnt\" { print \$1 }" "${ROOT}proc/mounts" | tail -n 1
}

find_partitions() {
	local dev="$1" type="$2" search=

	case "$type" in
	boot)
		search=bootable
		sfdisk -d "${ROOT}${dev#/}" | awk '/'$search'/ {print $1}'
		;;
	*)
		search=$(partition_id "$type")
		sfdisk -d "${ROOT}${dev#/}" | awk '/type='$search'/ {print $1}'
		;;
	esac
}

unmount_partitions() {
	local mnt="$1"
	umount $(awk '{print $2}' "${ROOT}proc/mounts" | grep -E "^$mnt(/|\$)" | sort -r)
}

init_progs() {
	local fs=

	for fs in $BOOTFS $ROOTFS $VARFS; do
		# we need load btrfs module early to avoid the error message:
		# 'failed to open /dev/btrfs-control'
		if ! grep -q -w "$fs" /proc/filesystems; then
			modprobe $fs
		fi
	done
	apk add --quiet sfdisk e2fsprogs $@
}

select_bootloader_pkg() {
	if [ "$BOOTLOADER" = none ]; then
		return
	fi
	local bootloader="$BOOTLOADER"
	echo "$bootloader"
}

find_nth_non_boot_parts() {
	local idx="$1" id="$2" disk=
	shift 2
	local disks="$@"

	for disk in $disks; do
		sfdisk -d "${ROOT}${disk#/}" | grep -v "bootable" | awk "/(Id|type)=$id/ { i++; if (i==$idx) print \$1 }"
	done
}

reset_var() {
	[ -z "$(find_mount_dev /var)" ] && return 0
	mkdir -pv /.var
	mv /var/* /.var/ 2>/dev/null
	umount /var && rm -rf /var && mv /.var /var && rm -rf /var/lost+found
}

setup_root() {
	local root_dev="$1"
	shift 2
	local disks="$@"
	local sysroot="${ROOT}${SYSROOT#/}"

	mkdir -p "$sysroot"
	if [ $(mountpoint -q "$sysroot" && echo 1 || echo 0) -eq 0 ]; then
		mount -t "$ROOTFS" "$root_dev" "$sysroot" || return 1
	fi

	install_mounted_root "$sysroot" "$disks" || return 1
	$MOCK swapoff -a
	unmount_partitions "$sysroot"
}

install_mounted_root() {
	local mnt="$(realpath "$1")"
	shift 1
	local disks="${@}"
	local rootdev=

	rootdev=$(find_mount_dev "$mnt")
	if [ -z "$rootdev" ]; then
		echo "'$mnt' does not seem to be a mount point" >&2
		return 1
	fi

	local apkflags="--initdb --progress --update-cache --clean-protected"
	local pkgs=

	pkgs="$pkgs $(cat /usr/sbin/apk-xorg | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-docker | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-pulseaudio | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-chromium | tr '\n' ' ')"

	local arch="$(apk --print-arch)"
	local repositories_file="$mnt"/etc/apk/repositories
	local keys_dir="$mnt"/etc/apk/keys

	init_chroot_mounts "$mnt"

	printf "\n\n"
	print_heading1 " Copying application"
	print_heading1 "----------------------"

	cp -vfrT /usr/sbin/hms/minikube "$mnt"/usr/bin/

	local export_dir="$mnt"/home/data/dist/
	makefile root:root 0644 "$mnt"/etc/profile.d/00-postinstall.sh <<-EOF
	EOF
	cat /usr/sbin/postinstall-hms.sh >>"$mnt"/etc/profile.d/00-postinstall.sh

	mkdir -pv $export_dir
	cp -vfr /usr/sbin/hms/app/* $export_dir
	for exe in \
		api \
		downloader \
		encoder \
		file; do
		cp -vfrT /usr/sbin/hms/$exe $export_dir
	done

	printf "\n\n"
	print_heading1 " Installing packages"
	print_heading1 "----------------------"

	apk add --keys-dir "$keys_dir" \
		--repositories-file "$repositories_file" \
		--no-script --root "$mnt" --arch "$arch" \
		$apkflags $pkgs
	local ret=$?
	cleanup_chroot_mounts "$mnt"
	return $ret
}

native_disk_install() {
	local root_part_type="$(partition_id linux)"
	local root_dev=
	init_progs $(select_bootloader_pkg) || return 1

	$MOCK mdev -s

	local index=1
	root_dev=$(find_nth_non_boot_parts $index "$root_part_type" $@)

	setup_root $root_dev $@
}

cat >>/etc/apk/repositories <<EOF
https://dl-cdn.alpinelinux.org/alpine/v3.21/main
https://dl-cdn.alpinelinux.org/alpine/v3.21/community
EOF
if rc-service --exists networking; then
	rc-service networking --quiet restart
fi

printf "y" | /usr/sbin/setup-alpine -e -f /usr/sbin/env-hms-answers.sh

reset_var
$MOCK swapoff -a

# stop all volume groups in use
$MOCK vgchange --ignorelockingfailure -a n >/dev/null 2>&1

diskdevs=
diskdevs="$diskdevs /dev/sda"

set -- $diskdevs
$MOCK dmesg -n1

native_disk_install $diskdevs

RC=$?
exit $RC
