profile_hms() {
    profile_standard
    profile_abbrev="hms"
    hostname="alpine"
    kernel_cmdline="console=tty0 console=ttyS0,115200"
    syslinux_serial="0 115200"
    arch="x86_64"
    kernel_flavors="lts"
    initfs_cmdline="modules=loop,squashfs,sd-mod,usb-storage quiet nouveau.config=NvGspRm=1"
    # initfs_cmdline="modules=loop,squashfs,sd-mod,usb-storage quiet"
    initfs_features="ata base bootchart cdrom ext4 mmc nvme raid scsi squashfs usb virtio"
    modloop_sign=yes
    grub_mod="all_video disk part_gpt part_msdos linux normal configfile search search_label efi_gop fat iso9660 cat echo ls test true help gzio"
    # x86_64
    grub_mod="$grub_mod multiboot2 efi_uga"
    initfs_features="$initfs_features nfit"

    apks="alpine-base apk-cron busybox chrony dhcpcd doas e2fsprogs kbd-bkeymaps network-extras openntpd openssl openssh tzdata wget tiny-cloud-alpine"
    rootfs_apks="busybox alpine-baselayout alpine-keys alpine-release apk-tools libc-utils"

    apks="$apks $(cat /tmp/build/apk-common | tr '\n' ' ') $(cat /tmp/build/apk-font | tr '\n' ' ') $(cat /tmp/build/apk-qt | tr '\n' ' ')"
    local _k _a
    for _k in $kernel_flavors; do
        apks="$apks linux-$_k"
        for _a in $kernel_addons; do
            apks="$apks $_a-$_k"
        done
    done
    apks="$apks linux-firmware"

    apkovl="aports/scripts/genapkovl-hms.sh"
}
