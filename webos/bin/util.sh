#!/bin/sh

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat >"$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

check_ready() {
	[ $# -lt 1 ] && return $? || continue

	local check_cmd=${1:-}
	local qualifier=${2:-1}
	local duration=${3:-600}
	local retry=0

	while [[ $(echo "$check_cmd" | sh | wc -l) -le $qualifier ]]; do
		if [[ $(echo "$check_cmd" | sh | wc -l) -ge $qualifier ]]; then break; fi
		if [ $retry -le $duration ]; then printf "\rWaiting service readiness (%s retried)..." $retry; else break; fi
		retry=$((retry + 1))
		sleep $CHECK_INTERVAL
	done

	sleep $CHECK_INTERVAL

	printf "\n"
	if [[ $(echo "$check_cmd" | sh | wc -l) -lt $qualifier ]]; then
		printf "%s\e[31m [failed]\e[0m\n" "$check_cmd"
		return $?
	fi
	printf "%s\e[32m [ok]\e[0m\n" "$check_cmd"
}
