FROM localhost:5000/gonode:latest AS build-iso
USER root
ENV GIT_SSL_NO_VERIFY=false
COPY \
    ./webos/*.sh \
    ./webos/apk-* \
    /tmp/build/
COPY ./dist/bin/ /tmp/bin/
COPY ./dist/app/ /tmp/app/
RUN apk update && apk add --no-cache \
        alpine-sdk \
        build-base \
        apk-tools \
        alpine-conf \
        ca-certificates \
        busybox \
        fakeroot \
        xorriso \
        squashfs-tools \
        sudo \
        mtools \
        dosfstools \
        grub-efi; \
    \
    apk add --no-cache \
        python3 \
        py3-pylspci \
        xorg-server; \
    \
    # asdasd
    addgroup root abuild; \
    mkdir -pv $HOME/.mkimage/; \
    mkdir -pv /export/iso/; \
    mkdir -pv $HOME/tmp/;
RUN \
    abuild-keygen -i -a -n; \
    ###########
    cd $HOME/.mkimage/; \
    git clone --depth=1 --branch=3.21-stable https://gitlab.alpinelinux.org/alpine/aports.git ./aports; \
    sudo apk update; \
    ###########
    export APORTS=$(realpath "$HOME/.mkimage/aports/"); \
    ###########
    cp -vr \
        /tmp/build/genapkovl-hms.sh \
        /tmp/build/mkimg.hms.sh \
        $APORTS/scripts/; \
    \
    sudo chmod u+x \
        $APORTS/scripts/genapkovl-hms.sh \
        $APORTS/scripts/mkimg.hms.sh; \
    \
    apk add --no-cache \
        pulseaudio \
        xorg-server \
        dbus \
        dbus-x11; \
    \
    $APORTS/scripts/mkimage.sh \
        --tag v3.21 \
        --outdir /export/iso/ \
        --arch x86_64 \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/main/ \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/community/ \
        --profile hms;

FROM scratch
COPY --from=build-iso /export/iso/ /
ENTRYPOINT ["/export/iso/"]
