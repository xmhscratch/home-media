#!/bin/sh

. "$(dirname $(realpath $0))/util.sh"

setup() {
	local mnt="$1"
	shift

	local export_dir=$(realpath "$mnt"/home/data/dist/)
	mkdir -pv "$export_dir"
	mkdir -pv "$export_dir"/bin/
	mkdir -pv "$export_dir"/channel/
	mkdir -pv "$export_dir"/node_modules/

	for exe in \
		api \
		downloader \
		encoder \
		file \
		tui \
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
	cp -fr /usr/sbin/hms/channel/* "$export_dir"/channel/
	cp -fr /usr/sbin/hms/node_modules/* "$export_dir"/node_modules/

	for exe in $(find /usr/sbin/hms/bin/* -type f | xargs basename -a); do
		install -o root -g root -m 0775 /usr/sbin/hms/bin/"$exe" "$export_dir"/bin/"$exe" ;;
	done

	for exe in $(find /usr/sbin/hms/sbin/* -type f | xargs basename -a); do
		install -o root -g root -m 0775 /usr/sbin/hms/sbin/"$exe" /usr/bin/hms/"$exe" ;;
	done

	# makefile root:wheel 0775 "$mnt"/usr/sbin/hms/tuid <<-EOF
	# socat -d -d UNIX-LISTEN:/run/tuid.sock,fork SYSTEM:'$export_dir/bin/tui'
	# EOF

	makefile root:wheel 0775 "$mnt"/etc/init.d/tuid <<-EOF
	#!/sbin/openrc-run

	command="/usr/sbin/hms/tuid"
	command_background=false
	name="tuid"

	depend() {
		need localmount
	}

	respawn_delay=2
	respawn_max=3
	respawn_period=60
	supervisor=supervise-daemon
	EOF

	rc_add tuid default
}

setup $1
