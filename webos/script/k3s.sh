#!/bin/sh

. "$(dirname $(realpath $0))/util.sh"

setup() {
	local mnt="$1"
	local usr=${2:-hms}
	shift 2

	local usr_home_dir=$(getent passwd "$(id -u $usr)" | cut -d: -f6)
	local ci_dir=/home/data/dist/ci/

	mkdir -pv "$mnt"/"$ci_dir"
	mkdir -pv "$mnt"/var/lib/rancher/k3s/agent/images/

	cp -vfr /usr/sbin/hms/*.yml "$mnt"/"$ci_dir"
	cp -vfrT \
		/usr/sbin/hms/k3s-airgap-images-amd64.tar.zst \
		"$mnt"/var/lib/rancher/k3s/agent/images/k3s-airgap-images-amd64.tar.zst
	cp -vfrT \
		/usr/sbin/hms/preload-images.tar.gz \
		"$mnt"/"$usr_home_dir"/preload-images.tar.gz
	gzip -d "$mnt"/"$usr_home_dir"/preload-images.tar.gz

	# openrc services and configuration files
	makefile root:wheel 0644 "$mnt"/etc/conf.d/k3s <<-EOF
	# k3s options
	export PATH=\$PATH:/usr/libexec/cni:/root/go/bin
	K3S_EXEC="server"
	K3S_OPTS="--disable=traefik --write-kubeconfig-mode=0644 --cluster-dns=10.43.0.10"
	rc_need="localmount netmount rpcbind"
	rc_after="net nfs"
	EOF

	makefile root:wheel 0775 "$mnt"/etc/init.d/k3s-proxy <<-EOF
	#!/sbin/openrc-run

	K3S_PROXY_LOGFILE="\${K3S_PROXY_LOGFILE:-/var/log/\${RC_SVCNAME}.log}"

	name="k3s-proxy"
	description="K3s apiserver proxy"
	command="/usr/bin/kubectl proxy"
	command_args="\${command_args:-} >>\${K3S_PROXY_LOGFILE}"
	command_background="no"
	pidfile="/var/run/\${RC_SVCNAME}.pid"

	start_pre() {
	    checkpath -f -m 0774 -o root:wheel "\${K3S_PROXY_LOGFILE}"
	}

	depend() {
	    after net k3s
	    provide k3s-proxy
	}

	respawn_delay=2
	respawn_max=3
	respawn_period=60
	supervisor=supervise-daemon
	EOF

	rc_add k3s dnsmasq default
}

setup $1 $2
