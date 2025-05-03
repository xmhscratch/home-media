#!/bin/sh

makefile() {
	OWNER="$1"
	PERMS="$2"
	FILENAME="$3"
	cat > "$FILENAME"
	chown "$OWNER" "$FILENAME"
	chmod "$PERMS" "$FILENAME"
}

rc_add() {
	mkdir -pv /etc/runlevels/"$2"
	if [ ! -f /etc/runlevels/"$2"/"$1" ]; then
		ln -sf /etc/init.d/"$1" /etc/runlevels/"$2"/"$1"
	fi
}

setup_user() {
	[[ -z "$(rc-status -q boot | grep polkit)" ]] && rc_add polkit boot

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

cfg_xorg() {
	[[ -z "$(rc-status -q boot | grep dbus)" ]] && rc_add dbus boot
	[[ -z "$(rc-status -q boot | grep seatd)" ]] && rc_add seatd boot

	makefile root:wheel 0644 /etc/X11/xorg.conf <<-EOF
	Section "ServerLayout"
	    Identifier     "X.org Configured"
	    Screen      0  "Screen0" 0 0
	    InputDevice    "Mouse0" "CorePointer"
	    InputDevice    "Keyboard0" "CoreKeyboard"
	EndSection

	Section "Files"
	    ModulePath   "/usr/lib/xorg/modules"
	    FontPath     "/usr/share/fonts/100dpi:unscaled"
        FontPath     "/usr/share/fonts/75dpi:unscaled"
		FontPath     "/usr/share/fonts/TTF"
        FontPath     "/usr/share/fonts/Type1"
	    FontPath     "/usr/share/fonts/cyrillic"
	    FontPath     "/usr/share/fonts/encodings"
	    FontPath     "/usr/share/fonts/misc"
	    FontPath     "/usr/share/fonts/noto"
	    FontPath     "/usr/share/fonts/opensans"
	EndSection

	Section "Module"
	    Load  "extmod"
	    Load  "glx"
	    Load  "dri"
	    Load  "dri2"
	    Load  "dbe"
	    Load  "record"
	EndSection

	Section "InputDevice"
	    Identifier  "Keyboard0"
	    Driver      "kbd"
	EndSection

	Section "InputDevice"
	    Identifier  "Mouse0"
	    Driver      "mouse"
	    Option      "Protocol" "auto"
	    Option      "Device" "/dev/input/mice"
	    Option      "ZAxisMapping" "4 5 6 7"
	EndSection
	EOF

	mkdir -pv /etc/X11/xorg.conf.d/

	makefile root:wheel 0644 /etc/X11/xorg.conf.d/10-monitor.conf <<-EOF
	Section "Monitor"
	    Identifier "Virtual-1"
	    UseModes "Modes-0"
	    Option "PreferredMode" "1366x768R"
	EndSection

	Section "Modes"
	    Identifier "Modes-0"
	    Modeline "4096x2160R"  567.00  4096 4144 4176 4256  2160 2163 2173 2222 +hsync -vsync
	    Modeline "2560x1600R"  268.50  2560 2608 2640 2720  1600 1603 1609 1646 +hsync -vsync
	    Modeline "1920x1200R"  154.00  1920 1968 2000 2080  1200 1203 1209 1235 +hsync -vsync
	    Modeline "1680x1050R"  119.00  1680 1728 1760 1840  1050 1053 1059 1080 +hsync -vsync
	    Modeline "1400x1050R"  101.00  1400 1448 1480 1560  1050 1053 1057 1080 +hsync -vsync
	    Modeline "1440x900R"   88.75  1440 1488 1520 1600  900 903 909 926 +hsync -vsync
	    Modeline "1368x768R"   72.25  1368 1416 1448 1528  768 771 781 790 +hsync -vsync
	    Modeline "1280x768R"   68.00  1280 1328 1360 1440  768 771 781 790 +hsync -vsync
	    Modeline "800x600R"   35.50  800 848 880 960  600 603 607 618 +hsync -vsync
	EndSection

	Section "Screen"
	    Identifier "Screen0"
	    Monitor "Virtual-1"
	    DefaultDepth 24
	    SubSection "Display"
	        Modes "4096x2160R"
	    EndSubSection
	    SubSection "Display"
	        Modes "2560x1600R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1920x1200R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1680x1050R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1400x1050R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1440x900R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1368x768R"
	    EndSubSection
	    SubSection "Display"
	        Modes "1280x768R"
	    EndSubSection
	    SubSection "Display"
	        Modes "800x600R"
	    EndSubSection
	EndSection
	EOF

	makefile root:wheel 0644 /etc/X11/xorg.conf.d/20-gpu.conf <<-EOF
	Section "Device"
	    Identifier  "Card0"
	    Driver      "modesetting"
	    BusID       "PCI:0:2:0"
	EndSection
	EOF
}

cfg_k8s_cluster() {
	[[ -z "$(rc-status -q default | grep docker)" ]] && rc_add docker default

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
		if [ -z "\$(lsmod | grep br_netfilter)" ]; then
			modprobe br_netfilter;
		fi
	}

	start() {
		ebegin "Starting minikube"
		\${MINIKUBE_BINARY} start --driver=none --cpus=max --memory=max \$command_args
		eend \$?
	}

	stop() {
		ebegin "Stopping minikube"
		\${MINIKUBE_BINARY} stop --all
		eend \$?
	}
	EOF

	[[ -z "$(rc-status -q default | grep minikube)" ]] && rc_add minikube default
	rc-service --quiet minikube start

	if [[ -z "$(which kubectl)" ]]; then
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
	if [[ -z "$(awk '/^::once:\/sbin\/openrc\ async -q/a' /etc/inittab)" ]] || [ $# -gt 0 ]; then
		sed -i "/::wait:\/sbin\/openrc default/a ::once:\/sbin\/openrc async -q" /etc/inittab
	fi

	local ntp_srvname=pool.ntp.org
	local ntp_srvip=$(getent ahosts $ntp_srvname | head -n 1 | cut -d"STREAM $ntp_srvname" -f1)
	sed -i "s@$ntp_srvname@$ntp_srvip@g" "$mnt"/etc/chrony/chrony.conf

	[[ -z "$(rc-status -q default | grep chrony)" ]] || rc-update del chronyd
	[[ -z "$(rc-status -q async | grep chrony)" ]] && rc_add chrony async
	[[ -z "$(rc-status -q default | grep pulseaudio)" ]] && rc_add pulseaudio default
}

postinstall() {
	[ ! -f "/home/.renovated" ] || return $?

	setup_user
	cfg_xorg
	cfg_misc
	# cfg_k8s_cluster

	touch /home/.renovated
}
postinstall
