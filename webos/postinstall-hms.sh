#!/bin/sh

setup_user() {
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
		tty \
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
		[[ ! -z "$(grep -E '^'$grp':' /etc/group)" ]] || continue
		for j in $(grep -E '^'$grp':' /etc/group | sed -e 's/^.*://' | tr ',' ' '); do
			if [ $j == "hms" ]; then in_group=1; fi
		done
		if [ $in_group -eq 0 ]; then
			adduser hms $grp
		fi
	done
}

# cfg_xorg() { }

cfg_k8s_cluster() {
	echo THAWED > /sys/fs/cgroup/freezer/freezer.state

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

	# if [[ -z "$(which kubectl)" ]]; then
	# 	minikube kubectl -- get po -A
	# 	alias kubectl="minikube kubectl --"
	# fi

	# local app_dir=$(realpath "/home/data/dist/")
	# kubectl apply -f "$app_dir"/ci/
}

postinstall() {
	[ ! -f "/home/.renovated" ] || return $?
	setup_user
	# cfg_xorg
	cfg_k8s_cluster
	rc-update add -s default async
	touch /home/.renovated
}

postinstall
