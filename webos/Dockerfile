FROM localhost:5000/gonode:latest AS build-kube

USER root
COPY ./genapkovl-hms.sh ./mkimg.hms.sh ./apk-deps.txt /tmp/build/
ENV KUBE_ARCH=x86_64
ENV KUBE_VERS=1.35.0
RUN apk update && apk add --no-cache \
        alpine-sdk \
        coreutils \
        wget \
        openssh \
        git \
        go \
        build-base \
        apk-tools \
        alpine-conf; \
    \
    # asdasd
    cd $HOME/; \
    git clone --depth=1 --branch=v$KUBE_VERS https://github.com/kubernetes/minikube.git ./minikube-$KUBE_VERS; \
    cd ./minikube-$KUBE_VERS; \
    make -j$(getconf _NPROCESSORS_ONLN) linux; \
    mv -vf ./out/minikube-linux-amd64 /home/minikube; \
    chmod +x /home/minikube; \
    ls -lha /home/;

FROM localhost:5000/gonode:latest AS build-iso

USER root
COPY ./genapkovl-hms.sh ./mkimg.hms.sh ./apk-deps.txt /tmp/build/
COPY --from=build-kube /home/minikube /tmp/build/
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
    apk add --no-cache \
        python3 \
        py3-pylspci \
        xorg-server; \
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
        /tmp/build/genapkovl-hms.sh \
        /tmp/build/mkimg.hms.sh \
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
COPY --from=build-iso /mnt/iso/ /
ENTRYPOINT ["/mnt/iso/"]
