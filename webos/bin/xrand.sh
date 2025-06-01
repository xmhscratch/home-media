#!/bin/sh

. "/usr/bin/hms/util.sh"

init() {
    # local output=${output:-'Virtual-1'}
    # local resmode=${resmode:-'Virtual-1'}

    # check_ready "ss -lx | grep -E \"/.X11-unix/X$(echo $DISPLAY | sed 's/[^0-9]//g')\"" 1 $CHECK_TIMEOUT
    # xrandr | awk '/\ connected/{print $1,$2}'
    # xrandr --output Virtual-1 --mode "1440x900R";
}

# init $1 $2
