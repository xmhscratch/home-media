#!/bin/sh

SBIN_DIR=/usr/sbin
. "$SBIN_DIR/script/util.sh"

ZFS_POOL_NAME=${ZFS_POOL_NAME:-"storage-repo"}
ZFS_ASHIFT=${ZFS_ASHIFT:-12}
ZFS_RAID_LEVEL=${ZFS_RAID_LEVEL:-} # stripe,mirror,raidz1,raidz2,raidz3
ZFS_MOUNTPOINT=${ZFS_MOUNTPOINT:-/home/storage/} # stripe,mirror,raidz1,raidz2,raidz3

find_available_disks() {
    local bypart=/dev/disk/by-partuuid
    local raid_disks=$(ls -l $bypart | grep -E './s[b-z]{2}[0-9]+' | awk -F' ' '{print $11}' | tr '\n' ' ')
}

setup() {
	local mnt="$1"
	shift

	makefile root:wheel 0664 "$mnt"/etc/udev/rules.d/01-hotplug.rules <<-EOF
    KERNEL=="sd[b-z]", SUBSYSTEM=="block", ACTION=="add", RUN+="/usr/bin/hotplug.sh"
	EOF

    local bypart=/dev/disk/by-partuuid
    local raid_disks=$(ls -l $bypart | grep -E './s[b-z]{2}[0-9]+' | awk -F' ' '{print $9}' | tr '\n' ' ')
    local disk_num=$($lsdisk | wc -l)

    local raid_level=$ZFS_RAID_LEVEL
    local disk_num=$($lsdisk | wc -l)

	case $disk_num in
	    0*) printf "RAID level \"%s\" is not supported!" $disk_num; die ;;
        1*) raid_level=$ZFS_RAID_LEVEL;;
        2*) raid_level=raidz1;;
        3*) raid_level=raidz2;;
        4*) raid_level=raidz3;;
	    *) raid_level=$ZFS_RAID_LEVEL;;
	esac

    # zpool create -R "/mnt" -O compression=lz4 -O sharenfs=on -O mountpoint=/home/storage/ -O atime=off -O relatime=off -O canmount=on -O exec=off -o cachefile=/etc/zfs/zpool.cache -o ashift=12 -m /home/test/ storage-repo /dev/sdb
    zpool_create_cmd="zpool create @ \
        -R "$mnt" \
        -m $ZFS_MOUNTPOINT \
        -O compression=lz4 \
        -O sharenfs=on \
        -O mountpoint=$ZFS_MOUNTPOINT \
        -O atime=off \
        -O relatime=off \
        -O canmount=on \
        -O exec=off \
        -o cachefile=/etc/zfs/zpool.cache \
        -o ashift=$ZFS_ASHIFT \
        $ZFS_POOL_NAME $raid_level $raid_disks"

    if [[ -z $(echo "\-n" | xargs -I @ sh -c "$zpool_create_cmd" >/dev/stdout 2>&1) ]]; then
        return 1
    fi
    echo "\ " | xargs -I @ sh -c "$zpool_create_cmd"

    rc_add zfs-import zfs-mount boot
}

setup $1
