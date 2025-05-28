#!/bin/sh

CHECK_TIMEOUT=${CHECK_TIMEOUT:-1000}
CHECK_INTERVAL=${CHECK_INTERVAL:-0.1}

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

json_display_format="{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{'\n'}{end}"
pod_ready_filter="PodReadyToStartContainers=True;Initialized=True;Ready=True;ContainersReady=True;PodScheduled=True;"

wait_system() {
	local ingress_nginx_check_cmd="kubectl get pods -n ingress-nginx -l app.kubernetes.io/instance=ingress-nginx,app.kubernetes.io/name=ingress-nginx -o jsonpath=\"$json_display_format\""
	check_ready "$ingress_nginx_check_cmd | grep Initialized=True" 3 $CHECK_TIMEOUT
	check_ready "$ingress_nginx_check_cmd | grep -E '^ingress\-nginx\-controller.+PodReadyToStartContainers=True;Initialized=True;Ready=True;ContainersReady=True;PodScheduled=True;'" 1 $CHECK_TIMEOUT

	local kube_system_check_cmd="kubectl get pods -n kube-system -o jsonpath=\"{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{'\n'}{end}\""
	check_ready "$kube_system_check_cmd | grep -E '$pod_ready_filter'" 4 $CHECK_TIMEOUT
}

wait_redis() {
	local redis_check_cmd="kubectl get pods -n redis -o jsonpath=\"$json_display_format\""
	check_ready "$redis_check_cmd | grep '$pod_ready_filter'" $(kubectl get pods -n redis --no-headers | wc -l) $CHECK_TIMEOUT
}

wait_logstash() {
	local logstash_check_cmd="kubectl get pods -n logstash -o jsonpath=\"$json_display_format\""
	check_ready "$logstash_check_cmd | grep '$pod_ready_filter'" $(kubectl get pods -n logstash --no-headers | wc -l) $CHECK_TIMEOUT
}

wait_hms() {
	local hms_check_cmd="kubectl get pods -n hms -o jsonpath=\"$json_display_format\""
	check_ready "$hms_check_cmd | grep '$pod_ready_filter'" $(kubectl get pods -n hms --no-headers | wc -l) $CHECK_TIMEOUT
}

wait_apiserver() {
	local usr_home_dir=$1
	shift

	# create dashboard access token
	local k3s_bearer_token=
	if [ ! -f "$usr_home_dir"/.k3s_token ]; then
		k3s_bearer_token=$(kubectl -n default create token kubernetes-dashboard)
		kubectl config set-credentials kubernetes-dashboard --token="$k3s_bearer_token"
		makefile root:wheel 0444 "$usr_home_dir"/.k3s_token <<-EOF
			$k3s_bearer_token
		EOF
	else
		k3s_bearer_token=$(cat "$usr_home_dir"/.k3s_token)
	fi

	local k3s_root_dir=/var/lib/rancher/k3s
	local k3s_tlscert_dir="$k3s_root_dir"/server/tls

	check_ready "curl -s --cert $k3s_tlscert_dir/client-admin.crt --key $k3s_tlscert_dir/client-admin.key --cacert $k3s_tlscert_dir/server-ca.crt -H \"Authorization: Bearer $k3s_bearer_token\" https://127.0.0.1:6443/readyz | xargs -I @ sh -c '[ \"@\" == \"ok\" ] && echo -e || false'" 1 $CHECK_TIMEOUT
	if ! rc-service k3s-proxy --quiet status; then
		rc-service k3s-proxy --ifnotstarted start
	fi
	check_ready "curl -s http://127.0.0.1:8001/readyz | xargs -I @ sh -c '[ \"@\" == \"ok\" ] && echo -e || false'" 1 $CHECK_TIMEOUT
}

init() {
	exportfs -afv

    udevadm control --reload
    udevadm trigger

	local usr_home_dir=$(getent passwd "$(id -u hms)" | cut -d: -f6)
	local ci_dir=/home/data/dist/ci/

	check_ready "ss -lx | grep -E '.*containerd\.sock\ '" 1 $CHECK_TIMEOUT
	check_ready "ss -lx | grep -E '.*kubelet\.sock\ '" 1 $CHECK_TIMEOUT

	if [ ! -f "$usr_home_dir"/.renovated ]; then
		mkdir -pv /home/storage/ && chmod 0755 /home/storage/
		mkdir -pv /home/data/db/ && chmod 0775 /home/data/db/
		mkdir -pv /home/data/channel/ && chmod 0664 /home/data/channel/

		printf "%s\n" "Importing preload container images..."
		ctr image import "$usr_home_dir"/preload-images.tar
		kubectl apply -f "$ci_dir"/dashboard-deploy.yml
		wait_apiserver "$usr_home_dir"
		kubectl apply -f "$ci_dir"/ingress-nginx-deploy.yml
		wait_system
		kubectl apply -f "$ci_dir"/hms-deploy.yml
		wait_redis
		kubectl apply -f "$ci_dir"/redis-svc-ingress.yml
		wait_logstash
		kubectl apply -f "$ci_dir"/logstash-svc-ingress.yml
		wait_hms
		kubectl apply -f "$ci_dir"/hms-svc-ingress.yml

		touch "$usr_home_dir"/.renovated
	else
		wait_apiserver "$usr_home_dir"
		wait_system
		wait_redis
		wait_logstash
		wait_hms
	fi
}

init
