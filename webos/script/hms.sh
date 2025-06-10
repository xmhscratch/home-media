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

	makefile root:wheel 0775 "$mnt"/etc/init.d/tuidx <<-EOF
	#!/sbin/openrc-run

	TUIDX_LOGFILE="\${TUIDX_LOGFILE:-/var/log/\${RC_SVCNAME}.log}"
	TUIDX_ERRFILE="\${TUIDX_ERRFILE:-/var/log/\${RC_SVCNAME}.err}"

	name="tuidx"
	command="/usr/bin/socat"
	command_args="unix-listen:/run/tuidx.sock,fork,perm=0700 SYSTEM:'/usr/bin/hms/tuidx.pl'"
	command_background="yes"
	pidfile="/run/\${RC_SVCNAME}.pid"
	output_log="\${TUIDX_LOGFILE}"
	error_log="\${TUIDX_ERRFILE}"

	depend() {
	    need net
	    provide tuidx
	}

	start_pre() {
	    checkpath -f -m 0744 -o root:wheel "\${TUIDX_LOGFILE}"
	    checkpath -f -m 0744 -o root:wheel "\${TUIDX_ERRFILE}"
	}

	respawn_delay=2
	respawn_max=3
	respawn_period=60
	EOF

	rc_add tuidx default
}

setup $1
