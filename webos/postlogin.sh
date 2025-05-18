#!/bin/sh

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat >"$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

cfg_k3s_cluster() {
	local usr_home_dir=$1
	local k3s_dir=$2
	shift 2

	local retry=0
	while [[ $( ss -lx | grep containerd | wc -l ) -eq 0 ]]; do
		[ $retry -le 600 ] && echo "Waiting for containerd ..." || break
		retry=$((retry + 1))
		sleep 0.5
	done

	if [[ $( ss -lx | grep containerd | wc -l ) -eq 0 ]]; then
		echo "containerd is not ready!"
		return $?
	fi

	ctr image import "$usr_home_dir"/preload-images.tar

	kubectl apply -f "$k3s_dir"/ci/ingress-nginx-deploy.yml
	kubectl apply -f "$k3s_dir"/ci/dashboard-deploy.yml

	local token=$(kubectl -n default create token kubernetes-dashboard)
	kubectl config set-credentials kubernetes-dashboard --token="$token"

	makefile root:wheel 0444 "$usr_home_dir"/.k3s_token <<-EOF
	$token
	EOF
}

init() {
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
	exportfs -afv

	mkdir -pv /home/data/db/
	mkdir -pv /home/data/dist/
	mkdir -pv /home/data/dist/channel/

	local usr_home_dir=$(getent passwd "$(id -u hms)" | cut -d: -f6)
	local k3s_dir="$usr_home_dir"/.rancher/k3s
	local k3s_root_dir=/var/lib/rancher/k3s
	local k3s_tlscert_dir="$k3s_root_dir"/server/tls

	if [ ! -f "$usr_home_dir"/.k3s_token ]; then
		cfg_k3s_cluster "$usr_home_dir" "$k3s_dir"
	fi
	local k3s_bearer_token=$(cat "$usr_home_dir"/.k3s_token)

	local check_k3s_ready_cmd="curl \
		-s --max-time 10 \
		--cert \"$k3s_tlscert_dir\"/client-admin.crt \
		--key \"$k3s_tlscert_dir\"/client-admin.key \
		--cacert \"$k3s_tlscert_dir\"/server-ca.crt \
		-H \"Authorization: Bearer $k3s_bearer_token\" \
		https://127.0.0.1:6443/healthz"

	local retry=0
	while [[ $( echo $check_k3s_ready_cmd | sh ) != "ok" ]]; do
		[ $retry -le 600 ] && echo "Waiting responses from kubernetes apiserver ..." || break
		retry=$((retry + 1))
		sleep 0.5
	done

	if [[ $( echo $check_k3s_ready_cmd | sh ) != "ok" ]]; then
		echo "Kubernetes apiserver responses timeout!"
		return $?
	fi

	{ rc-service k3s-proxy --ifnotstarted --quiet start; } &

	local check_ingress_ready_cmd="kubectl get pods \
		-n ingress-nginx \
		-l app.kubernetes.io/component=controller,app.kubernetes.io/instance=ingress-nginx,app.kubernetes.io/name=ingress-nginx \
		-o jsonpath='{range .items[*]}{.status.containerStatuses[*].ready.true}{.metadata.name}{ "\\n"}{end}'"

	retry=0
	while [[ $( echo $check_ingress_ready_cmd | wc -l ) -eq 0 ]]; do
		[ $retry -le 600 ] && echo "Waiting for nginx ingress controller ..." || break
		retry=$((retry + 1))
		sleep 0.5
	done

	if [[ $( echo $check_ingress_ready_cmd | wc -l ) -eq 0 ]]; then
		echo "No responses from Nginx ingress controller!"
		return $?
	fi

	if [ ! -f "$usr_home_dir"/.renovated ]; then
		kubectl apply -k "$k3s_dir"/ci/
	    touch "$usr_home_dir"/.renovated
	fi
}

init
