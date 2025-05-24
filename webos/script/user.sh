#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

setup() {
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

setup $1
