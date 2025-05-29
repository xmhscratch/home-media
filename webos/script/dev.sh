#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

setup() {
	local mnt="$1"
	local dev_server="$2"
	local repo_name="home-media"
	shift

	local arch="$(apk --print-arch)"
	local repositories_file="$mnt"/etc/apk/repositories
	local keys_dir="$mnt"/etc/apk/keys

	sed -Ei "s@[\#]{0,}PermitRootLogin.+@PermitRootLogin yes@g" "$mnt"/etc/ssh/sshd_config
	sed -Ei "s@[\#]{0,}PasswordAuthentication.+@PasswordAuthentication yes@g" "$mnt"/etc/ssh/sshd_config
	sed -Ei "s@[\#]{0,}PermitEmptyPasswords.+@PermitEmptyPasswords yes@g" "$mnt"/etc/ssh/sshd_config
	sed -Ei "s@[\#]{0,}AllowTcpForwarding.+@AllowTcpForwarding yes@g" "$mnt"/etc/ssh/sshd_config

	cat >>"$mnt"/etc/fstab <<-EOF
		$dev_server:/$repo_name/	/home/dev/	nfs	addr=$dev_server,nfsvers=4,exec,nodev,noatime,nodiratime,soft,rsize=1048576,wsize=1048576,X-mount.mkdir=0775   0 0
	EOF

	echo -e "https://dl-cdn.alpinelinux.org/alpine/v3.21/main/" >> /etc/apk/repositories;
	echo -e "https://dl-cdn.alpinelinux.org/alpine/v3.21/community/" >> /etc/apk/repositories;

	apk update && apk add npm;
	apk add \
		--keys-dir "$keys_dir" \
		--repositories-file "$repositories_file" \
		--root "$mnt" --arch "$arch" \
		--progress --update-cache --clean-protected \
		alpine-sdk nodejs go npm jq docker docker-openrc \
		$([ ! -z "$VIRT" ] && [ "$VIRT" -eq 1 ] && echo "zfs-virt") \
	;

	npm install --global --prefix --verbose "$mnt/$(npm root -g)" @angular/cli @nestjs/cli;
	cat /usr/sbin/hms/node-modules.txt | xargs npm cache add --cache "$mnt/$(npm config get cache)" --verbose;

	rc_add docker default
}

setup $1 $2
