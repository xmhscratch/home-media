#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

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
		case $exe in
			*) install -o root -g root -m 0755 /usr/sbin/hms/bin/"$exe" "$export_dir"/bin/"$exe" ;;
		esac
	done
}

setup $1
