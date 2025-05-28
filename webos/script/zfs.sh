#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

setup() {
	local mnt="$1"
	shift

	makefile root:wheel 0664 "$mnt"/etc/udev/rules.d/01-hotplug.rules <<-EOF
    KERNEL=="sd[b-z]", SUBSYSTEM=="block", ACTION=="add", RUN+="/usr/bin/hotplug.sh"
	EOF

    # stripe,mirror,raidz1,raidz2,raidz3
    local bypart=/dev/disk/by-partuuid
    local raid_level=
    local raid_disks=$(ls -l $bypart | grep -E '../../s[b-z]{2}[0-9]+' | awk -F' ' '{print $11}' | tr '\n' ' ')
    local disk_num=$($lsdisk | wc -l)

	case $disk_num in
	    0) die "RAID level \"$disk_num\" is not supported!" ;;
        1) raid_level=stripe;;
        2) raid_level=raidz1;;
        3) raid_level=raidz2;;
        4) raid_level=raidz3;;
	    *) raid_level=stripe;;
	esac

    # -o sharenfs=on \
    # -o share=name=storage,path=/home/storage,prot=nfsv4 \
    zpool create -f \
        -m /home/storage/ \
        -o mountpoint=/home/storage/ \
        -o ashift=12 \
        -o atime=off \
        -o relatime=off \
        -o compression=lz4 \
        -o canmount=on \
        -o exec=off \
        -o cachefile=/etc/zfs/zpool.cache \
        storage-repo $raid_level $raid_disks

    rc_add zfs-import zfs-mount boot
}

setup $1
