FROM alpine:3.21 AS build

USER root
COPY ./genapkovl-mkimgoverlay.sh ./mkimg.hms.sh /tmp/aports-scripts/
RUN apk update && apk add --no-cache \
        alpine-sdk \
        alpine-conf \
        syslinux \
        xorriso \
        squashfs-tools \
        grub \
        grub-efi \
        doas; \
    \
    # asdasd
    [ ! -e $(getent group "abuild") ] || addgroup abuild; \
    [ ! -e $(getent passwd "build") ] || adduser build -G abuild; \
    echo "permit nopass :abuild" > /etc/doas.d/doas.conf; \
    mkdir -p /home/build/.mkimage/; \
    mkdir -p /home/build/iso/; \
    chown -R build:abuild /home/build/;

USER build
RUN \
    printf "\n" | abuild-keygen -a; \
    # asdasd
    cd $HOME/.mkimage/; \
    git clone --depth=1 --branch=3.21-stable https://gitlab.alpinelinux.org/alpine/aports.git ./aports; \
    export APORTS=$(realpath "$HOME/.mkimage/aports/"); \
    export TMPDIR=/tmp; \
    echo $APORTS; \
    echo $TMPDIR; \
    cd $APORTS; \
    cp -vr /tmp/aports-scripts/genapkovl-mkimgoverlay.sh /tmp/aports-scripts/mkimg.hms.sh $APORTS/scripts/; \
    doas chmod +x $APORTS/scripts/genapkovl-mkimgoverlay.sh $APORTS/scripts/mkimg.hms.sh; \
    ls -lha $APORTS/scripts/; \
    # asdasd
    doas apk update; \
    $APORTS/scripts/mkimage.sh \
        --tag v3.21 \
        --outdir $HOME/iso \
        --arch x86_64 \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/main/ \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/community/ \
        --profile hms;

FROM scratch
COPY --from=build /home/build/iso/ /
ENTRYPOINT ["/home/build/iso/"]
