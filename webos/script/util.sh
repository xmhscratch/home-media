#!/bin/sh

if [ -t 1 ]; then
	COLCYAN="\e[36m"
	COLWHITE="\e[97m"
	COLRESET="\e[0m"
else
	COLCYAN=""
	COLWHITE=""
	COLRESET=""
fi

export COLCYAN
export COLWHITE
export COLRESET

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
