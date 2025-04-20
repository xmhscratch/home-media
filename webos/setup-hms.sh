#!/bin/sh

if [ -z /etc/init.d/ ]; then
    mkdir -pv /etc/init.d/
fi

if [ -z /etc/init.d/minikube ]; then
    cat > /etc/init.d/minikube << EOF
#!/sbin/openrc-run

name="Minikube Cluster"
description="Minikube is local Kubernetes, focusing on making it easy to learn and develop for Kubernetes."

command="${MINIKUBE_BINARY:-/usr/bin/minikube}"
command_args="${command_args:-${MINIKUBE_OPTS:-}}"
command_background="yes"

MINIKUBE_LOGFILE="${MINIKUBE_LOGFILE:-/var/log/${RC_SVCNAME}.log}"
MINIKUBE_LOGPROXY_OPTS="$MINIKUBE_LOGPROXY_OPTS -m"

export \
    MINIKUBE_LOGPROXY_CHMOD="${MINIKUBE_LOGPROXY_CHMOD:-0644}" \
    MINIKUBE_LOGPROXY_LOG_DIRECTORY="${MINIKUBE_LOGPROXY_LOG_DIRECTORY:-/var/log}" \
    MINIKUBE_LOGPROXY_ROTATION_SIZE="${MINIKUBE_LOGPROXY_ROTATION_SIZE:-104857600}" \
    MINIKUBE_LOGPROXY_ROTATION_TIME="${MINIKUBE_LOGPROXY_ROTATION_TIME:-86400}" \
    MINIKUBE_LOGPROXY_ROTATION_SUFFIX="${MINIKUBE_LOGPROXY_ROTATION_SUFFIX:-.%Y%m%d%H%M%S}" \
    MINIKUBE_LOGPROXY_ROTATED_FILES="${MINIKUBE_LOGPROXY_ROTATED_FILES:-5}"

pidfile="${MINIKUBE_PIDFILE:-/var/run/user/$(id -u)/$RC_SVCNAME.pid}"
output_logger="log_proxy $MINIKUBE_LOGPROXY_OPTS $MINIKUBE_LOGFILE"
rc_ulimit="${MINIKUBE_ULIMIT:--c unlimited -n 1048576 -u unlimited}"
retry="${MINIKUBE_RETRY:-TERM/60/KILL/10}"

depend() {
    need localmount
    after firewall docker
}

start_pre() {
    checkpath -f -m 0644 -o root:docker "$MINIKUBE_LOGFILE"
}

start_post() {
    if [ -z $(which kubectl) ]; then
        "$command" kubectl -- get po -A
        alias kubectl="minikube kubectl --"
    fi
    "$command" addons enable ingress
    "$command" addons enable ingress-dns
}

start() {
    ebegin "Starting minikube"
    "$command" start --driver=none --cpus=max --memory=max $command_args
    eend $?
}

stop() {
    ebegin "Stopping minikube"
    ${MINIKUBE_BINARY} stop --all
    eend $?
}
EOF
fi
