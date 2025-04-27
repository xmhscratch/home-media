#!/bin/sh

PREFIX=/usr
LIBDIR=$PREFIX/lib

. "$LIBDIR/libalpine.sh"
. "$LIBDIR/dasd-functions.sh"

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat > "$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

rc_add() {
	mkdir -p /etc/runlevels/"$2"
	ln -sf /etc/init.d/"$1" /etc/runlevels/"$2"/"$1"
}

###########
###########
mkdir -pv /etc/conf.d/
mkdir -pv /etc/profile.d/
mkdir -pv /etc/X11/xinit/
mkdir -pv /etc/sudoers.d/
mkdir -pv /etc/init.d/
mkdir -pv /usr/share/udhcpc/
###########
# linuxfb, wayland, eglfs, xcb, wayland-egl, minimalegl, minimal, offscreen, vkkhrdisplay, vnc
makefile root:root 0644 /etc/profile.d/00-default.sh <<EOF
export DISPLAY=:20 
export QT_QPA_PLATFORM=xcb
export XDG_SESSION_TYPE=x11
export WLR_BACKENDS=x11
export XDG_RUNTIME_DIR=/var/run/user/\$(id -u)
export XDG_VTNR=1
[ -z \$XDG_RUNTIME_DIR ] || mkdir -p \$XDG_RUNTIME_DIR;
if [[ -z \$XDG_RUNTIME_DIR ]]; then
    chmod 700 \$XDG_RUNTIME_DIR;
    chown \$(id -un):\$(id -gn) \$XDG_RUNTIME_DIR;
fi
export \$(dbus-launch)
EOF

makefile root:root 0644 /etc/X11/xinit/xserverrc <<EOF
#!/bin/sh
exec /usr/bin/X \$DISPLAY -config /etc/X11/xorg.conf -keeptty -nolisten tcp -novtswitch "\$@" vt\$XDG_VTNR
EOF

makefile root:root 0644 /etc/X11/Xwrapper.config <<EOF
needs_root_rights = no
allowed_users = anybody
EOF

makefile root:root 0644 /home/hms/.profile <<EOF
dbus-update-activation-environment DISPLAY XDG_CURRENT_DESKTOP XCURSOR_SIZE XCURSOR_THEME
if [ -n "\$DISPLAY" ] && [ "\$XDG_VTNR" -eq 1 ]; then
    startx ~/.xinitrc
fi
EOF

makefile root:root 0644 /home/hms/.xsession <<EOF
#!/bin/sh
xinit ~/.xinitrc
EOF

makefile root:root 0644 /home/hms/.xinitrc <<EOF
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
###########
makefile root:root 0644 /etc/init.d/minikube <<EOF
#!/sbin/openrc-run

name="Minikube Cluster"
description="Minikube is local Kubernetes, focusing on making it easy to learn and develop for Kubernetes."

command="\${MINIKUBE_BINARY:-/usr/bin/minikube}"
command_args="\${command_args:-\${MINIKUBE_OPTS:-}}"
command_background="yes"

MINIKUBE_LOGFILE="\${MINIKUBE_LOGFILE:-/var/log/\${RC_SVCNAME}.log}"
MINIKUBE_LOGPROXY_OPTS="\$MINIKUBE_LOGPROXY_OPTS -m"

export MINIKUBE_LOGPROXY_CHMOD="\${MINIKUBE_LOGPROXY_CHMOD:-0644}"
export MINIKUBE_LOGPROXY_LOG_DIRECTORY="\${MINIKUBE_LOGPROXY_LOG_DIRECTORY:-/var/log}"
export MINIKUBE_LOGPROXY_ROTATION_SIZE="\${MINIKUBE_LOGPROXY_ROTATION_SIZE:-104857600}"
export MINIKUBE_LOGPROXY_ROTATION_TIME="\${MINIKUBE_LOGPROXY_ROTATION_TIME:-86400}"
export MINIKUBE_LOGPROXY_ROTATION_SUFFIX="\${MINIKUBE_LOGPROXY_ROTATION_SUFFIX:-.%Y%m%d%H%M%S}"
export MINIKUBE_LOGPROXY_ROTATED_FILES="\${MINIKUBE_LOGPROXY_ROTATED_FILES:-5}"

pidfile="\${MINIKUBE_PIDFILE:-/var/run/user/\$(id -u)/\$RC_SVCNAME.pid}"
output_logger="log_proxy \$MINIKUBE_LOGPROXY_OPTS \$MINIKUBE_LOGFILE"
rc_ulimit="\${MINIKUBE_ULIMIT:--c unlimited -n 1048576 -u unlimited}"
retry="\${MINIKUBE_RETRY:-TERM/60/KILL/10}"

depend() {
    need localmount
    after firewall docker
}

start_pre() {
    checkpath -f -m 0644 -o root:docker "\$MINIKUBE_LOGFILE"
    if [ -z \$(lsmod | grep br_netfilter) ]; then
        modprobe br_netfilter;
    fi
}

start_post() {
    if [ -z \$(which kubectl) ]; then
        "\$command" kubectl -- get po -A
        alias kubectl="minikube kubectl --"
    fi
    "\$command" addons enable ingress
    "\$command" addons enable ingress-dns
}

start() {
    ebegin "Starting minikube"
    "\$command" start --driver=none --cpus=max --memory=max \$command_args
    eend \$?
}

stop() {
    ebegin "Stopping minikube"
    \${MINIKUBE_BINARY} stop --all
    eend \$?
}
EOF
###########
makefile root:root 0644 /etc/conf.d/pulseaudio <<EOF
EOF
cat /etc/conf.d/pulseaudio > /etc/conf.d/pulseaudio
cat >> /etc/conf.d/pulseaudio <<EOF
PULSEAUDIO_SHOULD_NOT_GO_SYSTEMWIDE=1
EOF

makefile root:root 0644 /etc/conf.d/dbus <<EOF
command_args="--system --nofork --nopidfile --syslog-only \${command_args:-}"
EOF
###########
Xorg :20 -depth 16 -nolisten tcp -configure
makefile root:root 0644 /etc/X11/xorg.conf <<EOF
EOF
cat ~/xorg.conf.new > /etc/X11/xorg.conf
###########
passwd -d hms

makefile root:root 0440 /etc/sudoers.d/hms <<EOF
hms ALL=(ALL) NOPASSWD: ALL
EOF

makefile root:root 0644 /usr/sbin/autologin <<EOF
#!/bin/sh
exec login -f hms
EOF

chmod +x /usr/sbin/autologin
sed -i 's@:respawn:/sbin/getty@:respawn:/sbin/getty -n -l /usr/sbin/autologin@g' /etc/inittab

for i in \
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
    input \
    tape \
    video \
    netdev \
    kvm \
    games \
    www-data \
    users \
    messagebus \
    docker \
    polkitd \
    seat \
    avahi \
    pulse-access \
    pulse \
; do
    in_group=0
    for j in $(grep -E '^'$1':' /etc/group | sed -e 's/^.*://' | tr ',' ' '); do
        if [ $j == "hms" ]; then in_group=1; fi
    done
    if [ $in_group -eq 0 ]; then adduser hms $i; fi
done

###########
rc_add dbus boot
rc_add seatd boot
rc_add chrony boot
rc_add pulseaudio boot
rc_add polkit default
rc_add docker default
###########
