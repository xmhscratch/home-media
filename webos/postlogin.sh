#!/bin/sh

CHECK_TIMEOUT=${CHECK_TIMEOUT:-600}

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
		if [[ $(echo "$check_cmd" | sh | wc -l) -eq $qualifier ]]; then break; fi
		if [ $retry -le $duration ]; then printf "\rWaiting service readiness (%s retried)..." $retry; else break; fi
		retry=$((retry + 1))
		sleep 0.5
	done

	sleep 0.5

	printf "\n"
	if [[ $(echo "$check_cmd" | sh | wc -l) -lt $qualifier ]]; then
		printf "%s\e[31m [failed]\e[0m\n" "$check_cmd"
		return $?
	fi
	printf "%s\e[32m [ok]\e[0m\n" "$check_cmd"
}

cfg_dashboard() {
	local k3s_dir=$1
	shift
	kubectl apply -f "$k3s_dir"/ci/dashboard-deploy.yml
}

cfg_bearer_token() {
	local usr_home_dir=$1
	shift

	local k3s_bearer_token=
	if [ -f "$usr_home_dir"/.k3s_token ]; then
		k3s_bearer_token=$(cat "$usr_home_dir"/.k3s_token)
	else
		k3s_bearer_token=$(kubectl -n default create token kubernetes-dashboard)
		kubectl config set-credentials kubernetes-dashboard --token="$k3s_bearer_token" >/dev/null
		makefile root:wheel 0444 "$usr_home_dir"/.k3s_token <<-EOF
		$k3s_bearer_token
		EOF
	fi
	echo $k3s_bearer_token
}

cfg_ingress_nginx() {
	local k3s_dir=$1
	shift
	kubectl apply -f "$k3s_dir"/ci/ingress-nginx-deploy.yml
}

cfg_hms() {
	local k3s_dir=$1
	shift
	kubectl apply -f "$k3s_dir"/ci/hms-deploy.yml
}

wait_ingress_nginx() {
	local base_check_cmd="kubectl get pods -n ingress-nginx -l app.kubernetes.io/instance=ingress-nginx,app.kubernetes.io/name=ingress-nginx -o jsonpath=\"{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{'\n'}{end}\""
	check_ready "$base_check_cmd | grep Initialized=True" 3 $CHECK_TIMEOUT
	check_ready "$base_check_cmd | grep -E '^ingress\-nginx\-controller.*Ready\=True'" 1 $CHECK_TIMEOUT
}

wait_hms() {
	local base_check_cmd="kubectl get pods -n hms -o jsonpath=\"{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{'\n'}{end}\""
	check_ready "$base_check_cmd | grep 'PodReadyToStartContainers=True;Initialized=True;Ready=True;ContainersReady=True;PodScheduled=True;'" 6 $CHECK_TIMEOUT
}

wait_apiserver() {
	local usr_home_dir=$1
	shift

	local k3s_bearer_token=$(cfg_bearer_token "$usr_home_dir")
	local k3s_root_dir=/var/lib/rancher/k3s
	local k3s_tlscert_dir="$k3s_root_dir"/server/tls
	local check_apiserver_ready_cmd="curl -s --cert $k3s_tlscert_dir/client-admin.crt --key $k3s_tlscert_dir/client-admin.key --cacert $k3s_tlscert_dir/server-ca.crt -H \"Authorization: Bearer $k3s_bearer_token\" https://127.0.0.1:6443/readyz | xargs -I @ sh -c '[ \"@\" == \"ok\" ] && echo -e || false'"

	check_ready "$check_apiserver_ready_cmd" 1 $CHECK_TIMEOUT

	{ rc-service k3s-proxy --ifnotstarted --quiet start; } &
	check_ready "curl -s http://127.0.0.1:8001/readyz | xargs -I @ sh -c '[ \"@\" == \"ok\" ] && echo -e || false'" 1 $CHECK_TIMEOUT
}

init() {
	for mod in \
		autofs4 \
		configs \
		nf_conntrack \
		br_netfilter; do
		[[ -z "$(lsmod | grep \"$mod\")" ]] && modprobe "$mod"
	done

	sysctl -p /etc/sysctl.conf
	exportfs -afv

	local usr_home_dir=$(getent passwd "$(id -u hms)" | cut -d: -f6)
	local k3s_dir="$usr_home_dir"/.rancher/k3s

	check_ready "ss -lx | grep -E '.*containerd\.sock\ '" 1 $CHECK_TIMEOUT
	if [ ! -f "$usr_home_dir"/.renovated ]; then
		mkdir -pv /home/data/db/

		printf "%s\n" "Importing preload container images..."
		ctr image import "$usr_home_dir"/preload-images.tar
		cfg_dashboard "$k3s_dir"
		wait_apiserver "$usr_home_dir"
		cfg_ingress_nginx "$k3s_dir"
		wait_ingress_nginx
		cfg_hms "$k3s_dir"
		wait_hms

		touch "$usr_home_dir"/.renovated
	else
		wait_apiserver "$usr_home_dir"
		wait_ingress_nginx
		wait_hms
	fi
}

init
