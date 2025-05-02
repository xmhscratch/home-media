#!/bin/sh

# Example answer file for setup-alpine script
# If you don't want to use a certain option, then comment it out

# Use US layout with US variant
KEYMAPOPTS="us us"

# Set hostname
HOSTNAMEOPTS=hms

# Set device manager to udev
DEVDOPTS=udev

# Contents of /etc/network/interfaces
INTERFACESOPTS="auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
hostname $HOSTNAMEOPTS
"

# Search domain of example.com, Cloudflare public nameserver
DNSOPTS="-d localhost 1.1.1.1"

# Set timezone to UTC
TIMEZONEOPTS="UTC"

# set http/ftp proxy
PROXYOPTS=none

# Add first mirror (CDN)
APKREPOSOPTS="-c -1"

# Create admin user
USEROPTS="-a -u -g wheel,audio,input,video,netdev hms"

# Install Openssh
SSHDOPTS=openssh

# Use chrony
NTPOPTS="chrony"

# Use /dev/sda as a sys disk
DISKOPTS="-m sys /dev/sda"

# Setup storage with label APKOVL for config storage
LBUOPTS="LABEL=APKOVL"
APKCACHEOPTS="/media/LABEL=APKOVL/cache"
