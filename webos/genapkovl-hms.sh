#!/bin/sh -e

HOSTNAME=hms
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
cat /tmp/build/apk-common >> "$tmp"/etc/apk/world

###########
echo "nameserver 1.1.1.1" >> "$tmp"/etc/resolv.conf
###########
mkdir -pv "$tmp"/usr/sbin/hms/

cp -vrf \
	/tmp/build/apk-* \
	/tmp/build/env-hms-answers.sh \
	/tmp/build/install-hms.sh \
	/tmp/build/postinstall-hms.sh \
	"$tmp"/usr/sbin/; \
\
cp -vrf /tmp/app/* "$tmp"/usr/sbin/hms/
cp -vf /tmp/bin/minikube "$tmp"/usr/sbin/hms/minikube

chmod +x \
	"$tmp"/usr/sbin/env-hms-answers.sh \
	"$tmp"/usr/sbin/install-hms.sh \
	"$tmp"/usr/sbin/postinstall-hms.sh; \
\
chmod +x "$tmp"/usr/sbin/hms/minikube
###########

rc_add devfs sysinit
rc_add dmesg sysinit
rc_add mdev sysinit
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

tar -c -C "$tmp" etc usr | gzip -9n > $HOSTNAME.apkovl.tar.gz