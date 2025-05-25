#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

setup() {
	local mnt="$1"
	shift

	if [ "$DEBUG" -eq 1 ]; then
		sed -Ei "s@[\#]{0,}PermitRootLogin.+@PermitRootLogin yes@g" "$mnt"/etc/ssh/sshd_config
		sed -Ei "s@[\#]{0,}PasswordAuthentication.+@PasswordAuthentication yes@g" "$mnt"/etc/ssh/sshd_config
		sed -Ei "s@[\#]{0,}PermitEmptyPasswords.+@PermitEmptyPasswords yes@g" "$mnt"/etc/ssh/sshd_config
		sed -Ei "s@[\#]{0,}AllowTcpForwarding.+@AllowTcpForwarding yes@g" "$mnt"/etc/ssh/sshd_config
	fi
}

setup $1
