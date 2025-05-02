#!/bin/sh

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat >"$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

rc_add() {
	mkdir -p /etc/runlevels/"$2"
	ln -sf /etc/init.d/"$1" /etc/runlevels/"$2"/"$1"
}

setup_user() {
	[[ -z $(rc-status -q boot | grep polkit) ]] && rc_add polkit boot

	for grp in \
		bin \
		daemon \
		sys \
		adm \
		disk \
		kmem \
		wheel \
		audio \
		cdrom \
		dialout \
		input \
		tape \
		video \
		netdev \
		kvm \
		games \
		www-data \
		users \
		messagebus \
		docker \
		polkitd \
		seat \
		avahi \
		pulse-access \
		pulse \
	; do
		in_group=0
		[[ ! -z $(grep -E '^'$grp':' /etc/group) ]] || continue
		for j in $(grep -E '^'$grp':' /etc/group | sed -e 's/^.*://' | tr ',' ' '); do
			if [ $j == "hms" ]; then in_group=1; fi
		done
		if [ $in_group -eq 0 ]; then
			adduser hms $grp
		fi
	done
}

cfg_xorg() {
	[[ -z $(rc-status -q boot | grep dbus) ]] && rc_add dbus boot
	[[ -z $(rc-status -q boot | grep seatd) ]] && rc_add seatd boot

	Xorg :20 -depth 16 -nolisten tcp -configure
	makefile root:wheel 0644 /etc/X11/xorg.conf <<-EOF
	EOF
	cat ~/xorg.conf.new > /etc/X11/xorg.conf
}

cfg_k8s_cluster() {
	[[ -z $(rc-status -q default | grep docker) ]] && rc_add docker default

	makefile root:wheel 0755 /etc/init.d/minikube <<-EOF
	#!/sbin/openrc-run

	name="Minikube Cluster"
	description="Minikube is local Kubernetes, focusing on making it easy to learn and develop for Kubernetes."

	command="\${MINIKUBE_BINARY:-/usr/bin/minikube}"
	command_args="\${command_args:-\${MINIKUBE_OPTS:-}}"
	command_background="yes"

	MINIKUBE_LOGFILE="\${MINIKUBE_LOGFILE:-/var/log/\${RC_SVCNAME}.log}"
	MINIKUBE_LOGPROXY_OPTS="\$MINIKUBE_LOGPROXY_OPTS -m"

	export MINIKUBE_LOGPROXY_CHMOD="\${MINIKUBE_LOGPROXY_CHMOD:-0644}"
	export MINIKUBE_LOGPROXY_LOG_DIRECTORY="\${MINIKUBE_LOGPROXY_LOG_DIRECTORY:-/var/log}"
	export MINIKUBE_LOGPROXY_ROTATION_SIZE="\${MINIKUBE_LOGPROXY_ROTATION_SIZE:-104857600}"
	export MINIKUBE_LOGPROXY_ROTATION_TIME="\${MINIKUBE_LOGPROXY_ROTATION_TIME:-86400}"
	export MINIKUBE_LOGPROXY_ROTATION_SUFFIX="\${MINIKUBE_LOGPROXY_ROTATION_SUFFIX:-.%Y%m%d%H%M%S}"
	export MINIKUBE_LOGPROXY_ROTATED_FILES="\${MINIKUBE_LOGPROXY_ROTATED_FILES:-5}"

	pidfile="\${MINIKUBE_PIDFILE:-/var/run/user/\$(id -u)/\$RC_SVCNAME.pid}"
	output_logger="log_proxy \$MINIKUBE_LOGPROXY_OPTS \$MINIKUBE_LOGFILE"
	rc_ulimit="\${MINIKUBE_ULIMIT:--c unlimited -n 1048576 -u unlimited}"
	retry="\${MINIKUBE_RETRY:-TERM/60/KILL/10}"

	depend() {
		need localmount
		after firewall docker
	}

	start_pre() {
		checkpath -f -m 0644 -o root:docker "\$MINIKUBE_LOGFILE"
		if [ -z \$(lsmod | grep br_netfilter) ]; then
			modprobe br_netfilter;
		fi
	}

	start() {
		ebegin "Starting minikube"
		"\$command" start --driver=none --cpus=max --memory=max \$command_args
		eend \$?
	}

	stop() {
		ebegin "Stopping minikube"
		\${MINIKUBE_BINARY} stop --all
		eend \$?
	}
	EOF

	[[ -z $(rc-status -q default | grep minikube) ]] && rc_add minikube default
	rc-service --quiet minikube start

	if [[ -z $(which kubectl) ]]; then
		minikube kubectl -- get po -A
		alias kubectl="minikube kubectl --"

		minikube addons enable ingress
		minikube addons enable ingress-dns
	fi

	local app_dir=$(realpath "/home/data/dist/")
	kubectl apply -f "$app_dir"/ci/
}

cfg_misc() {
	mkdir -pv /etc/runlevels/async
	rc-update add -s default async
	if [[ -z $(awk '/^::once:\/sbin\/openrc\ async -q/a' /etc/inittab) ]] || [ $# -gt 0 ]; then
		sed -i "/::wait:\/sbin\/openrc default/a ::once:\/sbin\/openrc async -q" /etc/inittab
	fi

	[[ -z $(rc-status -q default | grep chrony) ]] || rc-update del chronyd
	[[ -z $(rc-status -q async | grep chrony) ]] && rc_add chrony async
	[[ -z $(rc-status -q default | grep pulseaudio) ]] && rc_add pulseaudio default
}

postinstall() {
	[ -f "/home/.renovated" ] || return $?

	setup_user
	cfg_xorg
	cfg_misc
	# cfg_k8s_cluster

	touch /home/.renovated
}
postinstall
