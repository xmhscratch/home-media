#!/bin/sh

. "$(dirname $(realpath $0))/util.sh"

setup() {
	local mnt="$1"
	shift

	echo "net.ipv4.ip_forward = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.ipv6.conf.all.forwarding = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.bridge.bridge-nf-call-iptables = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.bridge.bridge-nf-call-ip6tables = 1" | tee -a "$mnt"/etc/sysctl.conf
	echo "net.netfilter.nf_conntrack_tcp_be_liberal = 1" | tee -a "$mnt"/etc/sysctl.conf

	sed -i 's/^\([^#]\+swap\)/#\1/' "$mnt"/etc/fstab
	cat >> "$mnt"/etc/fstab <<-EOF
	cgroup2	/sys/fs/cgroup			cgroup2	rw,nosuid,nodev,noexec,relatime,X-mount.mkdir=0775			0 0
	cgroup	/sys/fs/cgroup/freezer	cgroup	rw,nosuid,nodev,noexec,relatime,X-mount.mkdir=0775,freezer	0 0
	EOF

	makefile root:wheel 644 "$mnt"/etc/cgconfig.conf <<-EOF
	mount {
	    blkio = /sys/fs/cgroup/blkio;
	    memory = /sys/fs/cgroup/memory;
	    cpu = /sys/fs/cgroup/cpu;
	    cpuacct = /sys/fs/cgroup/cpuacct;
	    cpuset = /sys/fs/cgroup/cpuset;
	    devices = /sys/fs/cgroup/devices;
	    hugetlb = /sys/fs/cgroup/hugetlb;
	    freezer = /sys/fs/cgroup/freezer;
	    net_cls = /sys/fs/cgroup/net_cls;
	    pids = /sys/fs/cgroup/pids;
	}
	EOF
	sed -Ei "s@#rc_ulimit=\".+\"@rc_ulimit=\"-u unlimited -c unlimited\"@g" "$mnt"/etc/rc.conf
	sed -Ei "s@#rc_cgroup_mode=\".+\"@rc_cgroup_mode=\"unified\"@g" "$mnt"/etc/rc.conf
	sed -i "s@#rc_cgroup_controllers=\"\"@rc_cgroup_controllers=\"blkio cpu cpuacct cpuset io devices hugetlb memory net_cls net_prio pids\"@g" "$mnt"/etc/rc.conf
	sed -i "s@#rc_controller_cgroups=@rc_controller_cgroups=@g" "$mnt"/etc/rc.conf
	# sed -i "s@#rc_default_runlevel=@rc_default_runlevel=@g" "$mnt"/etc/rc.conf
	# sed -i "s@#rc_cgroup_memory=\"\"@rc_cgroup_memory=\"memory.memsw.limit_in_bytes 4194304\"@g" "$mnt"/etc/rc.conf

	for mod in $(cat /usr/sbin/hms/kernel-modules.txt); do
		if [[ -z "$(grep -qxF "$mod" "$mnt"/etc/modules)" ]]; then
			echo "$mod" | tee -a "$mnt"/etc/modules
		fi
	done

	local ntp_srvname=pool.ntp.org
	local ntp_srvip=$(getent ahosts $ntp_srvname | head -n 1 | awk -F' ' '{print $1}')
	sed -i "s@$ntp_srvname@$ntp_srvip@g" "$mnt"/etc/chrony/chrony.conf
	sed -i "s@FAST_STARTUP=no@FAST_STARTUP=yes@g" "$mnt"/etc/conf.d/chronyd

	# makefile root:wheel 0644 "$mnt"/etc/conf.d/pulseaudio <<-EOF
	# # Config file for /etc/init.d/pulseaudio
	# # \$Header: /var/cvsroot/gentoo-x86/media-sound/pulseaudio/files/pulseaudio.conf.d,v 1.6 2006/07/29 15:34:18 flameeyes Exp \$

	# # For more see "pulseaudio -h".

	# # Startup options
	# PA_OPTS="--log-target=syslog --disallow-module-loading=1"
	# PULSEAUDIO_SHOULD_NOT_GO_SYSTEMWIDE=1
	# EOF

	rc_add seatd cgroups boot
	# rc_add dbus pulseaudio default
	rc_add dbus default
}

setup $1
