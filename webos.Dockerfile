FROM localhost:5000/gonode-build:latest AS build-iso
USER root
ENV GIT_SSL_NO_VERIFY=false
ENV ARCH=x86_64
COPY \
    ./webos/*.sh \
    ./webos/apk-* \
    ./webos/packages.txt \
    /tmp/build/
COPY ./dist/bin/ /tmp/bin/
COPY ./dist/app/ /tmp/app/
COPY ./.apks/ /tmp/apks/
RUN \
    ###########
    addgroup root abuild; \
    mkdir -pv $HOME/.mkimage/; \
    mkdir -pv /export/iso/; \
    mkdir -pv $HOME/tmp/; \
    \
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
    archdir="/tmp/apks/$ARCH"; \
    mkdir -pv "$archdir"; \
    apk update; \
    \
    pkgs=$(cat /tmp/build/packages.txt | tr '\n' ' '); \
    apk fetch --link --recursive --output "$archdir" $pkgs; \
    if [[ $(find "$archdir"/*.apk -type f | wc -l) -gt 0 ]]; then \
        apk index \
            --allow-untrusted \
            --repository /tmp/apks/ \
            --rewrite-arch "$ARCH" \
            --index "$archdir"/APKINDEX.tar.gz \
            --output "$archdir"/APKINDEX.tar.gz \
            "$archdir"/*.apk; \
        abuild-sign "$archdir"/APKINDEX.tar.gz; \
    fi; \
    echo "/tmp/apks/" > /etc/apk/repositories; \
    apk update --allow-untrusted; \
    apk fix --upgrade --allow-untrusted alpine-keys; \
    \
    apk add \
        pulseaudio \
        xorg-server \
        dbus \
        dbus-x11; \
    \
    $APORTS/scripts/mkimage.sh \
        --tag v3.21 \
        --outdir /export/iso/ \
        --arch $ARCH \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/main/ \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/community/ \
        --profile hms;

FROM scratch
COPY --from=build-iso /export/iso/ /
ENTRYPOINT ["/export/iso/"]
