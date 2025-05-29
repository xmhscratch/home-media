#!/bin/sh

ZFS_POOL_NAME=${ZFS_POOL_NAME:-"storage-repo"}
ZFS_ASHIFT=${ZFS_ASHIFT:-12}

init() {
    set -e

    echo -e "Disk plugged in: $DEVNAME"

    if lsblk -n /dev/$DEVNAME | grep -q disk && blkid /dev/$DEVNAME; then
        printf "Ignore none empty disk: %s" $DEVNAME
        return 1
    fi

    if [ -z $(zpool list | grep $ZFS_POOL_NAME)]; then
        printf "ZFS pool \"%s\" is uninitialized!" $ZFS_POOL_NAME
        return 1
    fi

    local devname="$DEVNAME"
    local bypart=/dev/disk/by-partuuid
    local lsdisk=$(ls -l $bypart | grep $devname)
    local diskuuid=$(echo $lsdisk | awk -F' ' '{print $9}')
    local diskdev=$(realpath "$bypart/$(echo $lsdisk | awk -F' ' '{print $11}')")

    if [ -e $diskdev ]; then
        dd if=/dev/zero of=$bypart/$diskuuid bs=$((2**$ZFS_ASHIFT)) status=progress
        echo "label: gpt" | sfdisk --quiet "/dev/$devname"
        echo ',,L' | sfdisk --unit=S "/dev/$devname"
        echo -e 'y' | mkfs.ext4 "$diskdev"

        local raid_level=$(zpool status $ZFS_POOL_NAME | awk '/^\s*(mirror|raidz[1-3]?)/ {print $1}')
        case $raid_level in
            mirror*) raid_level=mirror;;
            raidz1*) raid_level=raidz1;;
            raidz2*) raid_level=raidz2;;
            raidz3*) raid_level=raidz3;;
            *) raid_level=;;
        esac

        zpool_add_cmd="zpool add @ -o ashift=$ZFS_ASHIFT $ZFS_POOL_NAME $raid_level \"$diskdev\""
        if [[ -z $(echo "\-n" | xargs -I @ sh -c "$zpool_add_cmd" >/dev/stdout 2>&1) ]]; then
            return 1
        fi
        echo "\ " | xargs -I @ sh -c "$zpool_add_cmd"
    fi
}

init
