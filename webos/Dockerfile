FROM alpine:3.21 AS build

USER root
COPY ./genapkovl-hms.sh ./mkimg.hms.sh /tmp/aports-scripts/
RUN apk update && apk add --no-cache \
        alpine-sdk \
        build-base \
        apk-tools \
        alpine-conf \
        busybox \
        fakeroot \
        xorriso \
        squashfs-tools \
        sudo \
        mtools \
        dosfstools \
        grub-efi; \
    \
    # asdasd
    addgroup root abuild; \
    mkdir -pv $HOME/.mkimage/; \
    mkdir -pv /mnt/iso/; \
    mkdir -pv $HOME/tmp/;

RUN \
    abuild-keygen -i -a -n; \
    ###########
    cd $HOME/.mkimage/; \
    git clone --depth=1 --branch=3.21-stable https://gitlab.alpinelinux.org/alpine/aports.git ./aports; \
    sudo apk update; \
    ###########
    export APORTS=$(realpath "$HOME/.mkimage/aports/"); \
    export TMPDIR=$HOME/tmp/; \
    echo $APORTS; \
    echo $TMPDIR; \
    ###########
    cp -vr \
        /tmp/aports-scripts/genapkovl-hms.sh \
        /tmp/aports-scripts/mkimg.hms.sh \
        $APORTS/scripts/; \
    \
    sudo chmod +x $APORTS/scripts/genapkovl-hms.sh $APORTS/scripts/mkimg.hms.sh; \
    ls -lha $APORTS/scripts/; \
    $APORTS/scripts/mkimage.sh \
        --tag v3.21 \
        --outdir /mnt/iso/ \
        --arch x86_64 \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/main/ \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/community/ \
        --profile hms;

FROM scratch
COPY --from=build /mnt/iso/ /
ENTRYPOINT ["/mnt/iso/"]
