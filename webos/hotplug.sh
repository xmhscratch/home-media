#!/bin/sh

init() {
    echo -e "Disk plugged in: $DEVNAME"

    if lsblk -n /dev/$DEVNAME | grep -q disk && blkid /dev/$DEVNAME; then
        return 1
    fi

    local ashift=12
    local devname="$DEVNAME"
    local bypart=/dev/disk/by-partuuid
    local lsdisk=$(ls -l $bypart | grep $diskdev)
    local diskuuid=$(echo $lsdisk | awk -F' ' '{print $9}')
    local diskdev=$(realpath "$bypart/$(echo $lsdisk | awk -F' ' '{print $11}')")

    dd if=/dev/zero of=$bypart/$diskuuid bs=$((2**$ashift)) status=progress
    echo "label: gpt" | sfdisk --quiet "/dev/$devname"
    echo ',,L' | sfdisk --unit=S "/dev/$devname"
    echo -e 'y' | mkfs.ext4 "$diskdev"

    zpool add -f -o ashift=$ashift storage-repo raidz1 "$diskdev"
}

init
