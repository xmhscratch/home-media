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

rc_add() {
	local srvs=${*: 0:-1}

	local sn=0
	if [ $(($# > 1)) -eq 1 ]; then sn=$(($# - 1)); fi
	shift $sn

	for srv in $(echo $srvs | rev | cut -d' ' -f 2- | rev); do
		[ ! -d "$mnt"/etc/runlevels/"$@"/ ] || mkdir -pv "$mnt"/etc/runlevels/"$@"/
		if [ ! -f "$mnt"/etc/runlevels/"$@"/"$srv" ]; then
			ln -sf /etc/init.d/"$srv" "$mnt"/etc/runlevels/"$@"/"$srv"
		fi
	done
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
	pkgs="$pkgs $(cat /usr/sbin/apk-font | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-kube | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-pulseaudio | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-chromium | tr '\n' ' ')"
	pkgs="$pkgs $(cat /usr/sbin/apk-auth | tr '\n' ' ')"

	local arch="$(apk --print-arch)"
	local repositories_file="$mnt"/etc/apk/repositories
	local keys_dir="$mnt"/etc/apk/keys

	init_chroot_mounts "$mnt"

	if [ "$DEBUG" -eq 1 ]; then
		sed -Ei "s@[\#]{0,}PermitRootLogin.+@PermitRootLogin yes@g" "$mnt"/etc/ssh/sshd_config
		sed -Ei "s@[\#]{0,}PasswordAuthentication.+@PasswordAuthentication yes@g" "$mnt"/etc/ssh/sshd_config
		sed -Ei "s@[\#]{0,}PermitEmptyPasswords.+@PermitEmptyPasswords yes@g" "$mnt"/etc/ssh/sshd_config
	fi

	printf "\n\n"
	print_heading1 " Install application"
	print_heading1 "----------------------"
	setup_app "$mnt"

	printf "\n\n"
	print_heading1 " Register user"
	print_heading1 "----------------------"
	setup_user "$mnt"

	printf "\n\n"
	print_heading1 " Setup Xorg"
	print_heading1 "----------------------"
	cfg_xorg "$mnt"

	printf "\n\n"
	print_heading1 " Finish installation"
	print_heading1 "----------------------"
	cfg_misc "$mnt"
	cfg_dnsmasq "$mnt"
	cfg_k3s "$mnt"

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

setup_app() {
	local mnt="$1"
	shift 1

	printf "\n\n"
	print_heading2 " Copy HMS files"
	print_heading2 "----------------------"

	local export_dir="$mnt"/home/data/dist/
	mkdir -pv "$export_dir"
	mkdir -pv "$export_dir"/node_modules/
	export_dir=$(realpath "$export_dir")

	for exe in \
		api \
		downloader \
		encoder \
		file \
	; do
		cp -vfrT /usr/sbin/hms/"$exe" "$export_dir"/"$exe"
	done
	for dir in \
		backend \
		frontend \
	; do
		mkdir -pv "$export_dir"/"$dir"
		cp -fr /usr/sbin/hms/"$dir"/* "$export_dir"/"$dir"
	done
	cp -fr /usr/sbin/hms/node_modules/* "$export_dir"/node_modules/

	for exe in $(find /usr/sbin/hms/bin/* -type f | xargs basename -a); do
		case $exe in
			*) install -o root -g root -m 0755 /usr/sbin/hms/bin/"$exe" "$mnt"/usr/bin/"$exe" ;;
		esac
	done

	printf "\n\n"
	print_heading2 " Install system apks"
	print_heading2 "----------------------"

	apk add --keys-dir "$keys_dir" \
		--repositories-file "$repositories_file" \
		--root "$mnt" --arch "$arch" \
		$apkflags $pkgs
}

setup_user() {
	local mnt="$1"
	shift

	if [[ -z $( id "hms" &>/dev/null ) ]]; then
		adduser -D "hms" &>/dev/null
	fi

	passwd -d hms
	mkdir -pv "$mnt"/etc/sudoers.d/

	makefile root:root 0440 "$mnt"/etc/sudoers.d/hms <<-EOF
	hms ALL=(ALL) NOPASSWD: ALL
	Defaults env_keep += "DISPLAY QT_QPA_PLATFORM WLR_BACKENDS XDG_SESSION_TYPE XDG_VTNR XDG_RUNTIME_DIR DBUS_SESSION_BUS_ADDRESS DBUS_SESSION_BUS_PID XDG_CURRENT_DESKTOP XCURSOR_SIZE XCURSOR_THEME"
	EOF

	makefile root:root 0755 "$mnt"/usr/sbin/autologin <<-EOF
	#!/bin/sh
	exec login -f hms
	EOF
	chmod +x "$mnt"/usr/sbin/autologin

	if [[ -z $(awk '/^:respawn:\/sbin\/getty -n -l \/usr\/sbin\/autologin/' "$mnt"/etc/inittab) ]] || [ $# -gt 0 ]; then
		sed -i 's@:respawn:/sbin/getty@:respawn:/sbin/getty -n -l /usr/sbin/autologin@g' "$mnt"/etc/inittab
	fi

	for grp in \
		bin \
		daemon \
		sys \
		adm \
		disk \
		kmem \
		wheel \
		audio \
		cdrom \
		dialout \
		tty \
		input \
		tape \
		video \
		netdev \
		kvm \
		games \
		www-data \
		users \
		messagebus \
		polkitd \
		seat \
		avahi \
		pulse-access \
		pulse \
	; do
		in_group=0
		[[ ! -z "$(grep -E '^'$grp':' /etc/group)" ]] || continue
		for j in $(grep -E '^'$grp':' /etc/group | sed -e 's/^.*://' | tr ',' ' '); do
			if [ $j == "hms" ]; then in_group=1; fi
		done
		if [ $in_group -eq 0 ]; then
			adduser hms $grp
		fi
	done
}

cfg_xorg() {
	local mnt="$1"
	shift

	mkdir -pv "$mnt"/etc/conf.d/
	mkdir -pv "$mnt"/etc/X11/xinit/
	mkdir -pv "$mnt"/etc/X11/xorg.conf.d/

	makefile root:wheel 0644 "$mnt"/etc/conf.d/dbus <<-EOF
	command_args="--system --nofork --nopidfile --syslog-only \${command_args:-}"
	EOF

	# linuxfb, wayland, eglfs, xcb, wayland-egl, minimalegl, minimal, offscreen, vkkhrdisplay, vnc
	makefile root:wheel 0755 "$mnt"/etc/profile.d/01-default.sh <<-EOF
	#!/bin/sh
	export DISPLAY=:20 
	export QT_QPA_PLATFORM=xcb
	export WLR_BACKENDS=x11
	export XDG_SESSION_TYPE=x11
	export XDG_VTNR=1
	export XDG_RUNTIME_DIR=/var/run/user/\$(id -u)
	if [ ! -d \$XDG_RUNTIME_DIR ]; then
	    sudo mkdir -pv \$XDG_RUNTIME_DIR;

	    sudo chmod 0700 \$XDG_RUNTIME_DIR;
	    sudo chown \$(id -un):\$(id -gn) \$XDG_RUNTIME_DIR;
	fi
	export \$(dbus-launch)
	export PATH=\$PATH:/usr/libexec/cni:/root/go/bin
	EOF

	sed -Ei "s@Defaults secure_path=\"(.+)\"@Defaults secure_path=\"\1:/usr/libexec/cni:/root/go/bin\"@g" "$mnt"/etc/sudoers

	makefile root:wheel 0755 "$mnt"/etc/X11/xinit/xserverrc <<-EOF
	#!/bin/sh
	exec /usr/bin/X \$DISPLAY -config /etc/X11/xorg.conf -nolisten tcp -novtswitch "\$@" vt\$XDG_VTNR
	EOF

	makefile root:wheel 0644 "$mnt"/etc/X11/Xwrapper.config <<-EOF
	needs_root_rights=yes
	allowed_users=anybody
	EOF

	makefile root:wheel 0644 "$mnt"/etc/X11/xorg.conf <<-EOF
	Section "ServerLayout"
	    Identifier     "X.org Configured"
	    Screen      0  "Screen0" 0 0
	    InputDevice    "Mouse0" "CorePointer"
	    InputDevice    "Keyboard0" "CoreKeyboard"
	EndSection

	Section "Files"
	    ModulePath   "/usr/lib/xorg/modules"
	    FontPath     "/usr/share/fonts/100dpi:unscaled"
        FontPath     "/usr/share/fonts/75dpi:unscaled"
	    FontPath     "/usr/share/fonts/TTF"
        FontPath     "/usr/share/fonts/Type1"
	    FontPath     "/usr/share/fonts/cyrillic"
	    FontPath     "/usr/share/fonts/encodings"
	    FontPath     "/usr/share/fonts/misc"
	    FontPath     "/usr/share/fonts/noto"
	    FontPath     "/usr/share/fonts/opensans"
	EndSection

	Section "Module"
	    Load  "extmod"
	    Load  "glx"
	    Load  "dri"
	    Load  "dri2"
	    Load  "dbe"
	    Load  "record"
	EndSection

	Section "InputDevice"
	    Identifier  "Keyboard0"
	    Driver      "kbd"
	EndSection

	Section "InputDevice"
	    Identifier  "Mouse0"
	    Driver      "mouse"
	    Option      "Protocol" "auto"
	    Option      "Device" "/dev/input/mice"
	    Option      "ZAxisMapping" "4 5 6 7"
	EndSection
	EOF

	makefile root:wheel 0644 "$mnt"/etc/X11/xorg.conf.d/10-monitor.conf <<-EOF
	Section "Monitor"
	    Identifier "Virtual-1"
	    UseModes "Modes-0"
	    Option "PreferredMode" "1366x768R"
	EndSection

	Section "Modes"
	    Identifier "Modes-0"
	    Modeline "4096x2160R"  567.00  4096 4144 4176 4256  2160 2163 2173 2222 +hsync -vsync
	    Modeline "2560x1600R"  268.50  2560 2608 2640 2720  1600 1603 1609 1646 +hsync -vsync
	    Modeline "1920x1200R"  154.00  1920 1968 2000 2080  1200 1203 1209 1235 +hsync -vsync
	    Modeline "1680x1050R"  119.00  1680 1728 1760 1840  1050 1053 1059 1080 +hsync -vsync
	    Modeline "1400x1050R"  101.00  1400 1448 1480 1560  1050 1053 1057 1080 +hsync -vsync
	    Modeline "1440x900R"   88.75  1440 1488 1520 1600  900 903 909 926 +hsync -vsync
	    Modeline "1368x768R"   72.25  1368 1416 1448 1528  768 771 781 790 +hsync -vsync
	    Modeline "1280x768R"   68.00  1280 1328 1360 1440  768 771 781 790 +hsync -vsync
	    Modeline "800x600R"   35.50  800 848 880 960  600 603 607 618 +hsync -vsync
	EndSection

	Section "Screen"
	    Identifier "Screen0"
	    Monitor "Virtual-1"
	    DefaultDepth 24
	    SubSection "Display"
	        Modes "4096x2160R"
	    EndSubSection
	    SubSection "Display"
	        Modes "2560x1600R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1920x1200R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1680x1050R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1400x1050R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1440x900R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1368x768R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1280x768R"
	    EndSubSection
	    SubSection "Display"
	        Modes "800x600R"
	    EndSubSection
	EndSection
	EOF

	makefile root:wheel 0644 "$mnt"/etc/X11/xorg.conf.d/20-gpu.conf <<-EOF
	Section "Device"
	    Identifier  "Card0"
	    Driver      "modesetting"
	    BusID       "PCI:0:2:0"
	EndSection
	EOF

	local usr_home_dir=$(getent passwd "hms" | cut -d: -f6)

	makefile root:wheel 0755 "$mnt"/"$usr_home_dir"/.profile <<-EOF
	#!/bin/sh
	sudo /usr/bin/postlogin.sh;
	dbus-update-activation-environment DISPLAY QT_QPA_PLATFORM WLR_BACKENDS XDG_SESSION_TYPE XDG_VTNR XDG_RUNTIME_DIR XDG_CURRENT_DESKTOP XCURSOR_SIZE XCURSOR_THEME;
	if [ -n "\$DISPLAY" ] && [ "\$XDG_VTNR" -eq 1 ]; then
		xinit ~/.xinitrc -- \$DISPLAY
	fi
	EOF

	makefile root:wheel 0755 "$mnt"/"$usr_home_dir"/.xinitrc <<-EOF
	#!/bin/sh
	{ sleep 1; xrandr --output Virtual-1 --mode "1368x768R"; } &
	exec dbus-launch --sh-syntax --exit-with-session chromium \
		--window-size="1368,768" \
		--window-position="0,0" \
		--app="http://127.0.0.1:8001/api/v1/namespaces/default/services/https:kubernetes-dashboard:/proxy/" \
		--no-sandbox \
		--kiosk \
		--start-fullscreen \
		--enable-features=UseOzonePlatform \
		--ozone-platform=x11 \
		--enable-unsafe-swiftshader \
		--enable-features=Vulkan,webgpu;
	EOF
}

cfg_misc() {
	local mnt="$1"
	shift

	cat >> "$mnt"/etc/exports <<-EOF
	/home/      *(rw,no_subtree_check,no_root_squash,fsid=0,anonuid=0,anongid=0)
	EOF

	echo "net.ipv4.ip_forward = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.ipv6.conf.all.forwarding = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.bridge.bridge-nf-call-iptables = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.bridge.bridge-nf-call-ip6tables = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.netfilter.nf_conntrack_tcp_be_liberal = 1" | tee -a "$mnt"/etc/sysctl.conf

	sed -i 's/^\([^#]\+swap\)/#\1/' "$mnt"/etc/fstab
	cat >> "$mnt"/etc/fstab <<-EOF
	cgroup2	/sys/fs/cgroup			cgroup2	rw,nosuid,nodev,noexec,relatime,X-mount.mkdir=0775			0 0
	cgroup	/sys/fs/cgroup/freezer	cgroup	rw,nosuid,nodev,noexec,relatime,X-mount.mkdir=0775,freezer	0 0
	EOF

	makefile root:wheel 644 "$mnt"/etc/cgconfig.conf <<-EOF
	mount {
	    blkio = /sys/fs/cgroup/blkio;
	    memory = /sys/fs/cgroup/memory;
	    cpu = /sys/fs/cgroup/cpu;
	    cpuacct = /sys/fs/cgroup/cpuacct;
	    cpuset = /sys/fs/cgroup/cpuset;
	    devices = /sys/fs/cgroup/devices;
	    hugetlb = /sys/fs/cgroup/hugetlb;
	    freezer = /sys/fs/cgroup/freezer;
	    net_cls = /sys/fs/cgroup/net_cls;
	    pids = /sys/fs/cgroup/pids;
	}
	EOF
	sed -Ei "s@#rc_ulimit=\".+\"@rc_ulimit=\"-u unlimited -c unlimited\"@g" "$mnt"/etc/rc.conf
	sed -Ei "s@#rc_cgroup_mode=\".+\"@rc_cgroup_mode=\"unified\"@g" "$mnt"/etc/rc.conf
	sed -i "s@#rc_cgroup_controllers=\"\"@rc_cgroup_controllers=\"blkio cpu cpuacct cpuset io devices hugetlb memory net_cls net_prio pids\"@g" "$mnt"/etc/rc.conf
	sed -i "s@#rc_controller_cgroups=@rc_controller_cgroups=@g" "$mnt"/etc/rc.conf
	# sed -i "s@#rc_default_runlevel=@rc_default_runlevel=@g" "$mnt"/etc/rc.conf
	# sed -i "s@#rc_cgroup_memory=\"\"@rc_cgroup_memory=\"memory.memsw.limit_in_bytes 4194304\"@g" "$mnt"/etc/rc.conf

	for mod in \
		autofs4 \
		configs \
		nf_conntrack \
		br_netfilter \
	; do
		if [[ -z "$(grep -qxF "$mod" "$mnt"/etc/modules)" ]]; then
			echo "$mod" | tee -a "$mnt"/etc/modules
		fi
	done

	local ntp_srvname=pool.ntp.org
	local ntp_srvip=$(getent ahosts $ntp_srvname | head -n 1 | awk -F' ' '{print $1}')
	sed -i "s@$ntp_srvname@$ntp_srvip@g" "$mnt"/etc/chrony/chrony.conf
	sed -i "s@FAST_STARTUP=no@FAST_STARTUP=yes@g" "$mnt"/etc/conf.d/chronyd

	makefile root:wheel 0644 "$mnt"/etc/conf.d/pulseaudio <<-EOF
	# Config file for /etc/init.d/pulseaudio
	# \$Header: /var/cvsroot/gentoo-x86/media-sound/pulseaudio/files/pulseaudio.conf.d,v 1.6 2006/07/29 15:34:18 flameeyes Exp \$

	# For more see "pulseaudio -h".

	# Startup options
	PA_OPTS="--log-target=syslog --disallow-module-loading=1"
	PULSEAUDIO_SHOULD_NOT_GO_SYSTEMWIDE=1
	EOF

	rc_add seatd cgroups boot
	rc_add dbus nfs pulseaudio default
}

cfg_dnsmasq() {
	local mnt="$1"
	shift

	mv -vf "$mnt"/etc/udhcpc/udhcpc.conf "$mnt"/etc/udhcpc/udhcpc.conf.bak
	makefile root:wheel 0644 "$mnt"/etc/udhcpc/udhcpc.conf <<-EOF
	RESOLV_CONF=yes
	EOF

	sed -i '1i nameserver 127.0.0.1' "$mnt"/etc/resolv.conf

	mv -vf "$mnt"/etc/dnsmasq.conf "$mnt"/etc/dnsmasq.conf.bak
	makefile root:wheel 0644 "$mnt"/etc/dnsmasq.conf <<-EOF
	port=53
	strict-order
	no-resolv
	no-poll
	address=/.hms/127.0.0.1
	local-service
	listen-address=127.0.0.1
	bind-interfaces
	conf-dir=/etc/dnsmasq.d/,*.conf
	EOF
}

cfg_k3s() {
	local mnt="$1"
	shift

	local usr_home_dir=$(getent passwd "$(id -u hms)" | cut -d: -f6)
	local k3s_dir="$usr_home_dir"/.rancher/k3s

	mkdir -pv "$mnt"/"$k3s_dir"
	mkdir -pv "$mnt"/var/lib/rancher/k3s/agent/images/

	cp -vfr /usr/sbin/hms/ci/ "$mnt"/"$k3s_dir"/ci/
	cp -vfrT \
		/usr/sbin/hms/k3s-airgap-images-amd64.tar.zst \
		"$mnt"/var/lib/rancher/k3s/agent/images/k3s-airgap-images-amd64.tar.zst
	cp -vfrT \
		/usr/sbin/hms/preload-images.tar.gz \
		"$mnt"/"$usr_home_dir"/preload-images.tar.gz
	gzip -d "$mnt"/"$usr_home_dir"/preload-images.tar.gz

	cp -vfrT /usr/sbin/postlogin.sh "$mnt"/usr/bin/postlogin.sh

	makefile root:wheel 0644 "$mnt"/etc/conf.d/k3s <<-EOF
	# k3s options
	export PATH=\$PATH:/usr/libexec/cni:/root/go/bin
	K3S_EXEC="server"
	K3S_OPTS="--disable=traefik --write-kubeconfig-mode=0644 --cluster-dns=10.43.0.10"
	rc_need="localmount netmount"
	rc_after="net nfs"
	EOF

	makefile root:wheel 0775 "$mnt"/etc/init.d/k3s-proxy <<-EOF
	#!/sbin/openrc-run

	K3S_PROXY_LOGFILE="\${K3S_PROXY_LOGFILE:-/var/log/\${RC_SVCNAME}.log}"

	name="k3s-proxy"
	description="K3s apiserver proxy"

	command="/usr/bin/kubectl proxy"
	command_args=">>\${K3S_PROXY_LOGFILE} 2>&1"

	start_pre() {
	    checkpath -f -m 0774 -o root:wheel "\${K3S_PROXY_LOGFILE}"
	}

	depend() {
	    need localmount netmount
	    after net k3s
	}

	respawn_delay=2
	respawn_max=3
	respawn_period=60
	supervisor=supervise-daemon
	EOF

	rc_add k3s dnsmasq default
}

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
