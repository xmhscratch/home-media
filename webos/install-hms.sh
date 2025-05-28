#!/bin/sh

PREFIX=/usr
LIBDIR=$PREFIX/lib
SBIN_DIR=/usr/sbin

. "$LIBDIR/libalpine.sh"
. "$LIBDIR/dasd-functions.sh"
. "$SBIN_DIR/script/util.sh"

ROOTFS=${ROOTFS:-ext4}
BOOTFS=${BOOTFS:-ext4}
VARFS=${VARFS:-ext4}
DISKLABEL=${DISKLABEL:-dos}

# default location for mounted root
SYSROOT=${SYSROOT:-/mnt}

# machine arch
ARCH=$(apk --print-arch)
BOOTLOADER=${BOOTLOADER:-syslinux}

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

	init_chroot_mounts "$mnt"

    makefile root:root 0644 "$mnt"/etc/motd <<-EOF
	EOF

	printf "\n\n"
	print_heading1 " Install alpine packages"
	print_heading1 "----------------------"
	setup_app "$mnt"

	printf "\n\n"
	print_heading1 " Register user"
	print_heading1 "----------------------"
	$SBIN_DIR/script/user.sh "$mnt"

	printf "\n\n"
	print_heading1 " Setup system"
	print_heading1 "----------------------"
	$SBIN_DIR/script/system.sh "$mnt"
	$SBIN_DIR/script/xorg.sh "$mnt"
	$SBIN_DIR/script/nfs.sh "$mnt"
	# $SBIN_DIR/script/zfs.sh "$mnt"
	$SBIN_DIR/script/dnsmasq.sh "$mnt"
	$SBIN_DIR/script/k3s.sh "$mnt"

	printf "\n\n"
	print_heading1 " Copying HMS files"
	print_heading1 "----------------------"
	$SBIN_DIR/script/hms.sh "$mnt"

	cleanup_chroot_mounts "$mnt"
	return $ret
}

setup_app() {
	local mnt="$1"
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
	pkgs="$pkgs $(cat /usr/sbin/apk-font | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-kube | tr '\n' ' ')"
	# pkgs="$pkgs $(cat /usr/sbin/apk-pulseaudio | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-chromium | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-auth | tr '\n' ' ')"

	local arch="$(apk --print-arch)"
	local repositories_file="$mnt"/etc/apk/repositories
	local keys_dir="$mnt"/etc/apk/keys

	printf "\n\n"
	print_heading2 " Install system apks"
	print_heading2 "----------------------"

	apk add --keys-dir "$keys_dir" \
		--repositories-file "$repositories_file" \
		--root "$mnt" --arch "$arch" \
		$apkflags $pkgs

	if [ ! -z "$DEV" ]; then
		$SBIN_DIR/script/dev.sh "$mnt" $DEV
	fi
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

if rc-service --exists networking; then
	rc-service networking --quiet restart
fi

echo -e "y" | /usr/sbin/setup-alpine -e -f /usr/sbin/env-hms-answers.sh

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
