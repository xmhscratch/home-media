#!/bin/sh -e

HOSTNAME="$1"
if [ -z "$HOSTNAME" ]; then
	echo "usage: $0 hostname"
	exit 1
fi

cleanup() {
	rm -rf "$tmp"
}

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat > "$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

rc_add() {
	mkdir -p "$tmp"/etc/runlevels/"$2"
	ln -sf /etc/init.d/"$1" "$tmp"/etc/runlevels/"$2"/"$1"
}

tmp="$(mktemp -d)"
trap cleanup EXIT

mkdir -p "$tmp"/etc
makefile root:root 0644 "$tmp"/etc/hostname <<EOF
$HOSTNAME
EOF

mkdir -p "$tmp"/etc/network
makefile root:root 0644 "$tmp"/etc/network/interfaces <<EOF
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
EOF

mkdir -p "$tmp"/etc/apk
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
alpine-base
EOF

###########
cat /tmp/build/apk-deps.txt >> "$tmp"/etc/apk/world
###########
# wayland
# wayland-dbg
# wayland-protocols
# wayland-libs-egl
# wayland-libs-server
# wayland-libs-cursor
# weston
# weston-xwayland
# weston-backend-wayland
# lxqt-themes
# lxqt-about
# lxqt-admin
# lxqt-config
# lxqt-panel
# lxqt-runner
# lxqt-archiver
# lxqt-policykit
# lxqt-openssh-askpass
# qtermwidget
# qterminal
# openbox
###########
branch=edge

VERSION_ID=$([ ! -f /etc/os-release ] || awk -F= '$1=="VERSION_ID" {print $2}' /etc/os-release)
case $VERSION_ID in
*_alpha*|*_beta*) branch=edge;;
*.*.*) branch=v${VERSION_ID%.*};;
*) branch=v3.21;;
esac

cat > "$tmp"/etc/apk/repositories <<EOF
https://dl-cdn.alpinelinux.org/alpine/$branch/main
https://dl-cdn.alpinelinux.org/alpine/$branch/community
EOF
###########
mkdir -pv "$tmp/usr/bin/"
cp /tmp/build/minikube "$tmp"/usr/bin/minikube
chmod +x ""$tmp"/usr/bin/minikube"

mkdir -p "${tmp}/etc/init.d/"
makefile root:root 0755 "$tmp/etc/init.d/minikube" << EOF
#!/sbin/openrc-run
name="Minikube Cluster"
description="Minikube is local Kubernetes, focusing on making it easy to learn and develop for Kubernetes."

command="/usr/bin/minikube start"
command_args="--force --cpus=max --memory=max ${command_args:-}"
command_background="yes"
pidfile="/run/$RC_SVCNAME.pid"

depend() {
	need localmount
	after bootmisc
	after docker
}
EOF
###########
mkdir -pv "$tmp"/root/
mkdir -pv "$tmp"/etc/conf.d/
mkdir -pv "$tmp"/usr/share/lxqt/

Xorg :20 -depth 32 -nolisten tcp -configure
makefile root:root 0644 "$tmp"/etc/X11/xorg.conf <<EOF
EOF
cat /root/xorg.conf.new > "$tmp"/etc/X11/xorg.conf

makefile root:root 0644 "$tmp"/root/.xinitrc <<EOF
EOF
cat /etc/X11/xinit/xinitrc > "$tmp"/root/.xinitrc
cat >> "$tmp"/root/.xinitrc <<EOF
exec dbus-run-session startlxqt
exec dbus-run-session chromium --app=https://www.youtube.com/ --no-sandbox --kiosk --start-fullscreen --enable-features=UseOzonePlatform --ozone-platform=x11 --enable-unsafe-swiftshader
EOF

makefile root:root 0644 "$tmp"/etc/conf.d/dbus <<EOF
command_args="--system --nofork --nopidfile --syslog-only ${command_args:-}"
EOF

if [ -n $(awk '/^nameserver/ {printf "%s ",$2}' /etc/resolv.conf) ] || [ $# -gt 0 ]; then
    sed -i "/export dns/a dns=\"$(echo '1.1.1.1' | tr ',' ' ')\"" /usr/share/udhcpc/default.script
fi

makefile root:root 0644 "$tmp"/etc/X11/Xwrapper.config <<EOF
needs_root_rights = no
allowed_users = anybody
EOF

# linuxfb, wayland, eglfs, xcb, wayland-egl, minimalegl, minimal, offscreen, vkkhrdisplay, vnc
makefile root:root 0644 "$tmp"/etc/profile.d/00-startup.sh <<EOF
export XDG_SESSION_TYPE=x11
export DISPLAY=:20
export QT_QPA_PLATFORM="eglfs"
export XDG_RUNTIME_DIR=/run/user/$(id -u)
[ -z $XDG_RUNTIME_DIR ] || mkdir -pv $XDG_RUNTIME_DIR;
if [[ -z $XDG_RUNTIME_DIR ]]; then
	chmod 700 $XDG_RUNTIME_DIR;
	chown $(id -un):$(id -gn) $XDG_RUNTIME_DIR;
fi
# dbus-run-session xinit lxqt -- $DISPLAY
lxqt-session

adduser nopwd seat # For seatd
adduser nopwd audio # For audio hardware.
adduser nopwd dialout # For serial consoles.
adduser nopwd kvm # For hardware virtualisation capabilities.
adduser nopwd qemu # Ditto
adduser nopwd gnupg # Smart-cards.
adduser nopwd power # For powerctl (power management tool).
adduser nopwd video # For video devices. See below.
setup-devd udev
EOF

# exec dbus-run-session -- startlxqt
makefile root:root 0644 "$tmp"/usr/share/lxqt/session.conf <<EOF
[General]
leave_confirmation=false

[Environment]
GTK_CSD=0
GTK_OVERLAY_SCROLLING=0

[Mouse]
cursor_size=18
cursor_theme=whiteglass
acc_factor=20
acc_threshold=10
left_handed=false

[Keyboard]
delay=350
interval=20
beep=false

[Font]
antialias=true
hinting=true
dpi=96
EOF
###########
rc_add dbus boot
rc_add polkit default
rc_add networking default
rc_add docker default
rc_add minikube default
###########

rc_add devfs sysinit
rc_add dmesg sysinit
rc_add udev sysinit
rc_add hwdrivers sysinit
rc_add modloop sysinit

rc_add hwclock boot
rc_add modules boot
rc_add sysctl boot
rc_add hostname boot
rc_add bootmisc boot
rc_add syslog boot

rc_add mount-ro shutdown
rc_add killprocs shutdown
rc_add savecache shutdown

tar -c -C "$tmp" etc proc usr root | gzip -9n > $HOSTNAME.apkovl.tar.gz
