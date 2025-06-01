#!/bin/sh

. "$(dirname $(realpath $0))/util.sh"

setup() {
	local mnt="$1"
	shift

	sed -Ei "s@OPTS_RPC_NFSD=\".*\"@OPTS_RPC_NFSD=\"8 -V 4 -V 4.1\"@g" "$mnt"/etc/conf.d/nfs
	sed -i "s@#rc_need=\"!rpc.statd\"@rc_need=\"!rpc.statd\"@g" "$mnt"/etc/conf.d/nfsclient

	cat >> "$mnt"/etc/exports <<-EOF
	/home/      *(rw,no_subtree_check,no_root_squash,fsid=0,anonuid=0,anongid=0)
	EOF

	rc_add nfs default
}

setup $1
