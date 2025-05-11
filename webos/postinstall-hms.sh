#!/bin/sh

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat >"$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

cfg_k3s_proxy() {
	makefile root:wheel 0775 /etc/init.d/k3s-proxy <<-EOF
	#!/sbin/openrc-run

	name="K3s proxy"
	description="K3s Proxy service"

	command="\${K3S_PROXY_BINARY:-/usr/bin/kubectl}"
	command_args="proxy \${K3S_PROXY_OPTS:-}"
	command_background="yes"

	pidfile="\${K3S_PROXY_PIDFILE:-/var/run/\$RC_SVCNAME.pid}"
	retry="\${K3S_PROXY_RETRY:-TERM/60/KILL/10}"

	depend() {
		need localmount
		after net k3s
	}
	EOF

	[[ -z "$(rc-status -q default | grep k3s-proxy)" ]] && rc-update add k3s-proxy default
	rc-service --quiet --ifnotstarted k3s-proxy start
}

cfg_k3s_cluster() {
	local usr_home_dir=$1
	shift

	local k3s_dir="$usr_home_dir"/.rancher/k3s
	local dashboard_dir="$k3s_dir"/dashboard
	local ingress_nginx_dir="$k3s_dir"/ingress-nginx

	ctr image import "$usr_home_dir"/preload-images.tar
	kubectl apply \
		-f "$ingress_nginx_dir"/deploy.yml
	kubectl apply \
		-f "$dashboard_dir"/deploy.yml \
		-f "$dashboard_dir"/admin-user.yml \
		-f "$dashboard_dir"/admin-user-role.yml
	kubectl apply -k "$k3s_dir"/ci/
	kubectl -n kubernetes-dashboard create token admin-user > "$usr_home_dir"/token.txt
}

cfg_async() {
	# mkdir -pv /etc/runlevels/async/
	# rc-update add -s default async

	# if [[ -z "$(awk '/^::once:\/sbin\/openrc\ async -q/a' /etc/inittab)" ]] || [ $# -gt 0 ]; then
	# 	sed -i "/::wait:\/sbin\/openrc default/a ::once:\/sbin\/openrc async -q" /etc/inittab
	# fi

	# if [ -f "/etc/runlevels/default/chronyd" ]; then
	# 	unlink /etc/runlevels/default/chronyd
	# 	ln -sf /etc/init.d/chronyd /etc/runlevels/async/chronyd
	# fi
}

postinstall() {
	[ ! -f "/home/.renovated" ] || return $?

	for mod in \
		autofs4 \
		configs \
		nf_conntrack \
		br_netfilter \
	; do
		if [[ -z "$(lsmod | grep \"$mod\")" ]]; then
			modprobe "$mod";
		fi
	done
	sysctl -p /etc/sysctl.conf

	local usr_home_dir=$(getent passwd "$(id -u hms)" | cut -d: -f6)
	local k3s_root_dir=/var/lib/rancher/k3s
	local k3s_tlscert_dir="$k3s_root_dir"/server/tls
	local k3s_bearer_token="$usr_home_dir"/token.txt

	local check_k3s_ready_cmd="curl \
		-s --max-time 10 \
		--cert \"$k3s_tlscert_dir\"/client-admin.crt \
		--key \"$k3s_tlscert_dir\"/client-admin.key \
		--cacert \"$k3s_tlscert_dir\"/server-ca.crt \
		-H \"Authorization: Bearer $(cat $k3s_bearer_token)\" \
		https://127.0.0.1:6443/healthz"

	while [[ $( echo $check_k3s_ready_cmd | sh ) != "ok" ]]; do
		echo "Waiting responses from kubernetes apiserver ..."
		sleep 5
	done

	if [[ $( echo $check_k3s_ready_cmd | sh ) != "ok" ]]; then
		echo "Kubernetes apiserver responses timeout!"
		return $?
	fi

	cfg_async
	cfg_k3s_proxy "$usr_home_dir"
	cfg_k3s_cluster "$usr_home_dir"

	touch /home/.renovated
}

postinstall
