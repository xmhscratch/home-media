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
cat /tmp/build/apk-common >> "$tmp"/etc/apk/world
cat /tmp/build/apk-font >> "$tmp"/etc/apk/world
cat /tmp/build/apk-qt >> "$tmp"/etc/apk/world
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
mkdir -pv "$tmp"/usr/bin/
cp /tmp/build/minikube "$tmp"/usr/bin/minikube
chmod +x "$tmp"/usr/bin/minikube
###########
mkdir -pv "$tmp"/root/
mkdir -pv "$tmp"/etc/conf.d/

# linuxfb, wayland, eglfs, xcb, wayland-egl, minimalegl, minimal, offscreen, vkkhrdisplay, vnc
makefile root:root 0644 "$tmp"/etc/profile.d/00-default.sh <<EOF
export DISPLAY=:20
export QT_QPA_PLATFORM=xcb
export XDG_SESSION_TYPE=x11
export WLR_BACKENDS=x11
export XDG_RUNTIME_DIR=/var/run/user/$(id -u)
export XDG_CURRENT_DESKTOP=sway
export XDG_SESSION_DESKTOP=sway
export XDG_VTNR=1
export SWAYSOCK=/run/user/$(id -u)/sway-ipc.$(id -u).$(pgrep -x sway).sock
[ -z $XDG_RUNTIME_DIR ] || mkdir -p $XDG_RUNTIME_DIR;
if [[ -z $XDG_RUNTIME_DIR ]]; then
	chmod 700 $XDG_RUNTIME_DIR;
	chown $(id -un):$(id -gn) $XDG_RUNTIME_DIR;
fi
export $(dbus-launch)
EOF

Xorg :20 -depth 16 -nolisten tcp -configure
makefile root:root 0644 "$tmp"/etc/X11/xorg.conf <<EOF
EOF
cat /root/xorg.conf.new > "$tmp"/etc/X11/xorg.conf

makefile root:root 0644 "$tmp"/etc/X11/xinit/xserverrc <<EOF
#!/bin/sh
exec /usr/bin/X $DISPLAY -config /etc/X11/xorg.conf -keeptty -nolisten tcp -novtswitch "$@" vt$XDG_VTNR
EOF

makefile root:root 0644 "$tmp"/root/.profile <<EOF
#!/bin/sh
dbus-update-activation-environment DISPLAY XDG_CURRENT_DESKTOP XCURSOR_SIZE XCURSOR_THEME
if [ -n "$DISPLAY" ] && [ "$XDG_VTNR" = 1 ]; then
  startx ~/.xinitrc
fi
EOF

makefile root:root 0644 "$tmp"/root/.xinitrc <<EOF
#!/bin/sh
exec dbus-launch --sh-syntax --exit-with-session chromium \
	--app=https://www.youtube.com/ \
	--no-sandbox \
	--kiosk \
	--start-fullscreen \
	--enable-features=UseOzonePlatform \
	--ozone-platform=x11 \
	--enable-unsafe-swiftshader \
	--enable-features=Vulkan,webgpu;
EOF

makefile root:root 0644 "$tmp"/etc/conf.d/dbus <<EOF
command_args="--system --nofork --nopidfile --syslog-only ${command_args:-}"
EOF

if [ -n $(awk '/^nameserver/ {printf "%s ",$2}' /etc/resolv.conf) ] || [ $# -gt 0 ]; then
    sed -i "/export dns/a dns=\"$(echo '1.1.1.1' | tr ',' ' ')\"" /usr/share/udhcpc/default.script
fi

makefile root:root 0440 "$tmp"/etc/sudoers.d/root <<EOF
root ALL=(ALL) NOPASSWD: ALL
EOF

makefile root:root 0644 "$tmp"/usr/sbin/autologin <<EOF
#!/bin/sh
exec login -f root
EOF
chmod +x "$tmp"/usr/sbin/autologin
cat /etc/inittab > "$tmp"/etc/inittab
sed -i 's@:respawn:/sbin/getty@:respawn:/sbin/getty -n -l /usr/sbin/autologin@g' "$tmp"/etc/inittab

makefile root:root 0644 "$tmp"/etc/X11/Xwrapper.config <<EOF
needs_root_rights = no
allowed_users = anybody
EOF

makefile root:root 0644 "$tmp"/etc/X11/Xwrapper.config <<EOF
EOF
###########
rc_add dbus boot
rc_add seatd boot
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
