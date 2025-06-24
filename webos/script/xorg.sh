#!/bin/sh

. "$(dirname $(realpath $0))/util.sh"

setup() {
	local mnt="$1"
	local usr=${2:-hms}
	shift 2

	mkdir -pv "$mnt"/etc/conf.d/
	mkdir -pv "$mnt"/etc/X11/xinit/
	mkdir -pv "$mnt"/etc/X11/xorg.conf.d/

	makefile root:wheel 0644 "$mnt"/etc/conf.d/dbus <<-EOF
	DBUS_LOGFILE="\${DBUS_LOGFILE:-/var/log/dbus.log}"
	command_args="--system --nofork --nopidfile --syslog-only \${command_args:-} >>\${DBUS_LOGFILE}"
	EOF

	# linuxfb, wayland, eglfs, xcb, wayland-egl, minimalegl, minimal, offscreen, vkkhrdisplay, vnc
	makefile root:wheel 0755 "$mnt"/etc/profile.d/01-default.sh <<-EOF
	#!/bin/sh
	export DISPLAY=:20 
	export QT_QPA_PLATFORM=xcb
	export WLR_BACKENDS=x11
	export XDG_SESSION_TYPE=x11
	export XDG_VTNR=1
	export XDG_RUNTIME_DIR=/var/run/user/\$(id -u)
	if [ ! -d \$XDG_RUNTIME_DIR ]; then
	    sudo mkdir -pv \$XDG_RUNTIME_DIR;

	    sudo chmod 0700 \$XDG_RUNTIME_DIR;
	    sudo chown \$(id -un):\$(id -gn) \$XDG_RUNTIME_DIR;
	fi
	export \$(dbus-launch)
	export PATH=\$PATH:/usr/libexec/cni:/root/go/bin
	EOF

	sed -Ei "s@Defaults secure_path=\"(.+)\"@Defaults secure_path=\"\1:/usr/libexec/cni:/root/go/bin\"@g" "$mnt"/etc/sudoers

	makefile root:wheel 0755 "$mnt"/etc/X11/xinit/xserverrc <<-EOF
	#!/bin/sh
	exec /usr/bin/X \$DISPLAY -config /etc/X11/xorg.conf -nolisten tcp -novtswitch "\$@" vt\$XDG_VTNR
	EOF

	makefile root:wheel 0644 "$mnt"/etc/X11/Xwrapper.config <<-EOF
	needs_root_rights=yes
	allowed_users=anybody
	EOF

	makefile root:wheel 0644 "$mnt"/etc/X11/xorg.conf <<-EOF
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

	makefile root:wheel 0644 "$mnt"/etc/X11/xorg.conf.d/10-monitor.conf <<-EOF
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

	makefile root:wheel 0644 "$mnt"/etc/X11/xorg.conf.d/20-gpu.conf <<-EOF
	Section "Device"
	    Identifier  "Card0"
	    Driver      "modesetting"
	    BusID       "PCI:0:2:0"
	EndSection
	EOF

	local usr_home_dir=$(getent passwd "$usr" | cut -d: -f6)

	makefile root:wheel 0755 "$mnt"/"$usr_home_dir"/.profile <<-EOF
	#!/bin/sh
	sudo /usr/bin/hms/startup.sh;
	dbus-update-activation-environment DISPLAY QT_QPA_PLATFORM WLR_BACKENDS XDG_SESSION_TYPE XDG_VTNR XDG_RUNTIME_DIR XDG_CURRENT_DESKTOP XCURSOR_SIZE XCURSOR_THEME;
	if [ -n "\$DISPLAY" ] && [ "\$XDG_VTNR" -eq 1 ]; then
	    xinit ~/.xinitrc -- \$DISPLAY
	fi
	EOF

	# --app="http://127.0.0.1:8001/api/v1/namespaces/default/services/https:kubernetes-dashboard:/proxy/" \
	# --app="http://frontend.hms/" \
	makefile root:wheel 0755 "$mnt"/"$usr_home_dir"/.xinitrc <<-EOF
	#!/bin/sh
	{ sleep 1; xrandr --output Virtual-1 --mode "1440x900R"; } &
	exec dbus-launch --sh-syntax --exit-with-session chromium \
	    --window-size="1440,900" \
	    --window-position="0,0" \
		--app="http://127.0.0.1:8001/api/v1/namespaces/default/services/https:kubernetes-dashboard:/proxy/" \
	    --no-sandbox \
	    --kiosk \
	    --start-fullscreen \
	    --ozone-platform=x11 \
	    --enable-unsafe-swiftshader \
	    --enable-unsafe-webgpu \
	    --enable-features=UseOzonePlatform,Vulkan,VulkanFromANGLE,DefaultANGLEVulkan,webgpu;
	EOF
}

setup $1 $2
