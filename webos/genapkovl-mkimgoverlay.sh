#!/bin/sh -e

HOSTNAME="$1"
if [ -z "$HOSTNAME" ]; then
    echo "usage: $0 hostname"
    exit 1
fi

cleanup() {
    rm -rf "$tmp"
}

makefile() {
    OWNER="$1"
    PERMS="$2"
    FILENAME="$3"
    cat >"$FILENAME"
    chown "$OWNER" "$FILENAME"
    chmod "$PERMS" "$FILENAME"
}

rc_add() {
    mkdir -p "$tmp"/etc/runlevels/"$2"
    ln -sf /etc/init.d/"$1" "$tmp"/etc/runlevels/"$2"/"$1"
}

tmp="$(mktemp -d)"
trap cleanup EXIT

mkdir -p "$tmp"/etc
makefile root:root 0644 "$tmp"/etc/hostname <<EOF
$HOSTNAME
EOF

mkdir -p "$tmp"/etc/network
makefile root:root 0644 "$tmp"/etc/network/interfaces <<EOF
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
EOF

mkdir -p "$tmp"/etc/apk
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
alpine-base
EOF

###########
# repositories_file=/etc/apk/repositories
# keys_dir=/etc/apk/keys
# branch=edge

# VERSION_ID=$(awk -F= '$1=="VERSION_ID" {print $2}' /etc/os-release)
# case $VERSION_ID in
# *_alpha*|*_beta*) branch=edge;;
# *.*.*) branch=v${VERSION_ID%.*};;
# esac

cat > "$tmp"/etc/apk/repositories <<EOF
https://dl-cdn.alpinelinux.org/alpine/v3.21/main
https://dl-cdn.alpinelinux.org/alpine/v3.21/community
EOF
###########
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
xorg-server
xf86-input-libinput
xinit
eudev
EOF
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
mesa-dri-gallium
elogind
polkit-elogind
EOF
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
polkit-qt-dev
polkit-qt5
qt5-assistant
qt5-qdbusviewer
qt5-qt3d
qt5-qt3d-dev
qt5-qtbase
qt5-qtbase-dbg
qt5-qtbase-dev
qt5-qtbase-mysql
qt5-qtbase-odbc
qt5-qtbase-postgresql
qt5-qtbase-sqlite
qt5-qtbase-tds
qt5-qtbase-x11
qt5-qtcharts
qt5-qtcharts-dev
qt5-qtconnectivity
qt5-qtconnectivity-dev
qt5-qtdatavis3d
qt5-qtdatavis3d-dev
qt5-qtdeclarative
qt5-qtdeclarative-dbg
qt5-qtdeclarative-dev
qt5-qtfeedback
qt5-qtfeedback-dev
qt5-qtgamepad
qt5-qtgamepad-dev
qt5-qtgraphicaleffects
qt5-qtimageformats
qt5-qtkeychain
qt5-qtkeychain-lang
qt5-qtlocation
qt5-qtlocation-dev
qt5-qtlottie
qt5-qtlottie-dev
qt5-qtmultimedia
qt5-qtmultimedia-dev
qt5-qtnetworkauth
qt5-qtnetworkauth-dev
qt5-qtpim
qt5-qtpim-dev
qt5-qtpurchasing
qt5-qtpurchasing-dev
qt5-qtquick3d
qt5-qtquick3d-dev
qt5-qtquickcontrols
qt5-qtquickcontrols2
qt5-qtquickcontrols2-dev
qt5-qtquicktimeline
qt5-qtremoteobjects
qt5-qtremoteobjects-dev
qt5-qtscript
qt5-qtscript-dev
qt5-qtscxml
qt5-qtscxml-dev
qt5-qtsensors
qt5-qtsensors-dev
qt5-qtserialbus
qt5-qtserialbus-dev
qt5-qtserialport
qt5-qtserialport-dev
qt5-qtspeech
qt5-qtspeech-dev
qt5-qtsvg
qt5-qtsvg-dev
qt5-qtsystems
qt5-qtsystems-dev
qt5-qttools
qt5-qttools-dev
qt5-qttranslations
qt5-qtusb
qt5-qtusb-dev
qt5-qtvirtualkeyboard
qt5-qtvirtualkeyboard-dev
qt5-qtwayland
qt5-qtwayland-dev
qt5-qtwebchannel
qt5-qtwebchannel-dev
qt5-qtwebengine
qt5-qtwebengine-dev
qt5-qtwebglplugin
qt5-qtwebglplugin-dev
qt5-qtwebsockets
qt5-qtwebsockets-dev
qt5-qtwebsockets-libs
qt5-qtwebview
qt5-qtwebview-dev
qt5-qtx11extras
qt5-qtx11extras-dev
qt5-qtxmlpatterns
qt5-qtxmlpatterns-dev
EOF
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
mesa-vulkan-ati
mesa-vulkan-intel
mesa-vulkan-layers
mesa-vulkan-swrast
vulkan-headers
vulkan-loader
vulkan-loader-dbg
vulkan-loader-dev
vulkan-tools
EOF
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
chromium
chromium-chromedriver
chromium-dbg
chromium-lang
chromium-qt5
chromium-swiftshader
EOF
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
lxqt-build-tools
libqtxdg
libsysstat
liblxqt
libdbusmenu-lxqt
lxqt-globalkeys
lxqt-powermanagement
lxqt-session
lxqt-qtplugin
lxqt-notificationd
pm-utils
xdg-utils
EOF
makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
lximage-qt
obconf-qt
pavucontrol-qt
arandr
sddm
font-dejavu
dbus
dbus-x11
openbox
elogind
polkit-elogind
gvfs
udisks2
adwaita-qt
oxygen
EOF
# makefile root:root 0644 "$tmp"/etc/apk/world <<EOF
# lxqt-themes
# lxqt-about
# lxqt-admin
# lxqt-config
# lxqt-panel
# lxqt-runner
# lxqt-archiver
# lxqt-policykit
# lxqt-openssh-askpass
# lxqt-sudo
# qtermwidget
# qterminal
# openbox
# EOF
###########
makefile root:root 0644 "$tmp"/etc/modules <<EOF
EOF
###########

rc_add devfs sysinit
rc_add dmesg sysinit
rc_add mdev sysinit
rc_add hwdrivers sysinit
rc_add modloop sysinit

rc_add hwclock boot
rc_add modules boot
rc_add sysctl boot
rc_add hostname boot
rc_add bootmisc boot
rc_add syslog boot

rc_add mount-ro shutdown
rc_add killprocs shutdown
rc_add savecache shutdown

###########
rc_add elogind boot
rc_add polkit boot
rc_add xinit boot
rc_add dbus boot
###########

tar -c -C "$tmp" etc | gzip -9n >$HOSTNAME.apkovl.tar.gz
