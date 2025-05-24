#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

setup() {
	local mnt="$1"
	shift

	mv -vf "$mnt"/etc/udhcpc/udhcpc.conf "$mnt"/etc/udhcpc/udhcpc.conf.bak
	makefile root:wheel 0644 "$mnt"/etc/udhcpc/udhcpc.conf <<-EOF
	RESOLV_CONF=yes
	EOF

	sed -i '1i nameserver 127.0.0.1' "$mnt"/etc/resolv.conf

	mv -vf "$mnt"/etc/dnsmasq.conf "$mnt"/etc/dnsmasq.conf.bak
	makefile root:wheel 0644 "$mnt"/etc/dnsmasq.conf <<-EOF
	port=53
	strict-order
	no-resolv
	no-poll
	address=/.hms/127.0.0.1
	local-service
	listen-address=127.0.0.1
	bind-interfaces
	conf-dir=/etc/dnsmasq.d/,*.conf
	EOF
}

setup $1
