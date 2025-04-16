profile_hms() {
    profile_standard
    profile_abbrev="hms"
    hostname="alpine"
    kernel_cmdline="console=tty0 console=ttyS0,115200"
    syslinux_serial="0 115200"
    arch="aarch64 x86 x86_64"
    kernel_flavors="lts"
    initfs_cmdline="modules=loop,squashfs,sd-mod,usb-storage quiet"
    initfs_features="ata base bootchart cdrom ext4 mmc nvme raid scsi squashfs usb virtio"
    modloop_sign=yes
    grub_mod="all_video disk part_gpt part_msdos linux normal configfile search search_label efi_gop fat iso9660 cat echo ls test true help gzio"
    # x86_64
    grub_mod="$grub_mod multiboot2 efi_uga"
    initfs_features="$initfs_features nfit"

    apks="alpine-base apk-cron busybox chrony dhcpcd doas e2fsprogs kbd-bkeymaps network-extras openntpd openssl openssh tzdata wget tiny-cloud-alpine"
    rootfs_apks="busybox alpine-baselayout alpine-keys alpine-release apk-tools libc-utils"

    local _k _a
    for _k in $kernel_flavors; do
        apks="$apks linux-$_k"
        for _a in $kernel_addons; do
            apks="$apks $_a-$_k"
        done
    done

    apks="$apks
            xorg-server
            xf86-input-libinput
            xinit
            eudev
            mesa-dri-gallium
            elogind
            polkit-elogind
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
            mesa-vulkan-ati
            mesa-vulkan-intel
            mesa-vulkan-layers
            mesa-vulkan-swrast
            vulkan-headers
            vulkan-loader
            vulkan-loader-dbg
            vulkan-loader-dev
            vulkan-tools
            chromium
            chromium-chromedriver
            chromium-dbg
            chromium-lang
            chromium-qt5
            chromium-swiftshader
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
            lximage-qt
            obconf-qt
            pavucontrol-qt
            arandr
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
        "
    apkovl="aports/scripts/genapkovl-hms.sh"
}
