#!/bin/sh

# setup_user() { }

# cfg_xorg() { }

cfg_k8s_cluster() {
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

	# local usr_home_dir="$mnt"/home/hms
	# local dashboard_dir="$usr_home_dir"/.rancher/k3s/dashboard

	# ctr image import /home/hms/preload-images.tar
	# kubectl create -f "$dashboard_dir"/deploy.yaml -f "$dashboard_dir"/admin-user.yaml -f "$dashboard_dir"/admin-user-role.yaml
	# kubectl -n kubernetes-dashboard create token admin-user > "$dashboard_dir"/token.txt
	# kubectl proxy >/dev/null &

	# local app_dir=$(realpath "/home/data/dist/")
	# kubectl apply -f "$app_dir"/ci/
}

# cfg_async() {
# 	mkdir -pv "$mnt"/etc/runlevels/async/
	
# 	if [[ -z "$(awk '/^::once:\/sbin\/openrc\ async -q/a' "$mnt"/etc/inittab)" ]] || [ $# -gt 0 ]; then
# 		sed -i "/::wait:\/sbin\/openrc default/a ::once:\/sbin\/openrc async -q" "$mnt"/etc/inittab
# 	fi

# 	local ntp_srvname=pool.ntp.org
# 	local ntp_srvip=$(getent ahosts $ntp_srvname | head -n 1 | cut -d"STREAM $ntp_srvname" -f1)
# 	sed -i "s@$ntp_srvname@$ntp_srvip@g" "$mnt"/etc/chrony/chrony.conf

# 	if [ -f "$mnt"/etc/runlevels/default/chronyd ]; then
# 		unlink "$mnt"/etc/runlevels/default/chronyd
# 		ln -sf /etc/init.d/chronyd "$mnt"/etc/runlevels/async/chronyd
# 	fi

# 	for srv in \
# 		pulseaudio \
# 	; do
# 		if [ ! -f "$mnt"/etc/runlevels/default/"$srv" ]; then
# 			ln -sf /etc/init.d/"$srv" "$mnt"/etc/runlevels/default/"$srv"
# 		fi
# 	done
# }

postinstall() {
	[ ! -f "/home/.renovated" ] || return $?
	# setup_user
	# cfg_xorg
	cfg_k8s_cluster
	touch /home/.renovated
}

postinstall
