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
	mkdir -pv "$tmp"/etc/runlevels/"$2"
	ln -sf /etc/init.d/"$1" "$tmp"/etc/runlevels/"$2"/"$1"
}

tmp="$(mktemp -d)"
trap cleanup EXIT

mkdir -pv "$tmp"/etc
makefile root:root 0644 "$tmp"/etc/hostname <<EOF
$HOSTNAME
EOF

mkdir -pv "$tmp"/etc/network
makefile root:root 0644 "$tmp"/etc/network/interfaces <<EOF
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
EOF

makefile root:root 0644 "$tmp"/etc/modules <<EOF
$(cat /tmp/build/kernel-modules.txt)
EOF

mkdir -pv "$tmp"/etc/apk
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
alpine-base
EOF
cat /tmp/build/apk-common >> "$tmp"/etc/apk/world

###########
echo "nameserver 1.1.1.1" >> "$tmp"/etc/resolv.conf
###########
mkdir -pv "$tmp"/usr/sbin/hms/
mkdir -pv "$tmp"/usr/sbin/hms/bin/
mkdir -pv "$tmp"/usr/sbin/hms/sbin/
mkdir -pv "$tmp"/usr/sbin/hms/channel/
mkdir -pv "$tmp"/usr/sbin/hms/node_modules/
mkdir -pv "$tmp"/usr/sbin/script/

# files must be on both directories
cp -vf \
	/tmp/build/apk-* \
	"$tmp"/usr/sbin/ \
;
# move files
mv -vf \
	/tmp/build/env-hms-answers.sh \
	/tmp/build/install-hms.sh \
	"$tmp"/usr/sbin/ \
;
mv -vf /tmp/build/script/* "$tmp"/usr/sbin/script/
mv -vf /tmp/app/* "$tmp"/usr/sbin/hms/
mv -vf /tmp/bin/* "$tmp"/usr/sbin/hms/bin/
mv -vf /tmp/sbin/* "$tmp"/usr/sbin/hms/sbin/
mv -vf /tmp/channel/* "$tmp"/usr/sbin/hms/channel/
mv -vf /tmp/node_modules/* "$tmp"/usr/sbin/hms/node_modules/
mv -vf \
	/tmp/ci/*.yml \
	/tmp/node-modules.txt \
	/tmp/build/kernel-modules.txt \
	/tmp/docker/preload-images.tar.gz \
	/tmp/docker/k3s-airgap-images-amd64.tar.zst \
	"$tmp"/usr/sbin/hms/; \
\
chmod +x \
	"$tmp"/usr/sbin/hms/bin/* \
	"$tmp"/usr/sbin/hms/sbin/* \
	"$tmp"/usr/sbin/script/*.sh \
	"$tmp"/usr/sbin/env-hms-answers.sh \
	"$tmp"/usr/sbin/install-hms.sh \
;
###########
SBIN_DIR="$tmp"/usr/sbin/
$SBIN_DIR/script/user.sh "$tmp" root
makefile root:root 0644 "$tmp"/etc/inittab <<EOF
# /etc/inittab

::sysinit:/sbin/openrc sysinit
::sysinit:/sbin/openrc boot
::wait:/sbin/openrc default

# Set up a couple of getty's
tty1::respawn:/sbin/getty -n -l /usr/sbin/autologin 38400 tty1
tty2::respawn:/sbin/getty -n -l /usr/sbin/autologin 38400 tty2
tty3::respawn:/sbin/getty -n -l /usr/sbin/autologin 38400 tty3
tty4::respawn:/sbin/getty -n -l /usr/sbin/autologin 38400 tty4
tty5::respawn:/sbin/getty -n -l /usr/sbin/autologin 38400 tty5
tty6::respawn:/sbin/getty -n -l /usr/sbin/autologin 38400 tty6

# Put a getty on the serial port
#ttyS0::respawn:/sbin/getty -n -l /usr/sbin/autologin -L 115200 ttyS0 vt100

# Stuff to do for the 3-finger salute
::ctrlaltdel:/sbin/reboot

# Stuff to do before rebooting
::shutdown:/sbin/openrc shutdown
EOF
###########
mkdir -pv "$tmp"/root/
makefile root:root 0755 "$tmp"/etc/motd <<-EOF
EOF
makefile root:root 0755 "$tmp"/root/.profile <<-EOF
#!/bin/sh
. "/usr/lib/libalpine.sh"
ask_yesno "Choose 'y' to includes dev environment, 'n' for minimal system (y/n)"
case \$? in
    0) echo "/usr/sbin/install-hms.sh" | /bin/sh ;;
    *)
		ask "What's the IP address of devserver? [192.168.56.55]" 192.168.56.55
		echo "DEV=\$resp /usr/sbin/install-hms.sh" | /bin/sh
	;;
esac
EOF
###########
mkdir -pv "$tmp"/etc/init.d/
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

tar -c -C "$tmp" etc usr root | gzip -9n > $HOSTNAME.apkovl.tar.gz
