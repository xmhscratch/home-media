FROM localhost:5000/gonode-build:latest AS build-kube
USER root
ENV KUBE_ARCH=x86_64
ENV KUBE_VERS=1.35.0
ENV GIT_SSL_NO_VERIFY=false
RUN \
    # build minikube
    cd /home/; \
    git clone --depth=1 --branch=v$KUBE_VERS https://github.com/kubernetes/minikube.git ./minikube-$KUBE_VERS; \
    cd ./minikube-$KUBE_VERS; \
    make -j$(getconf _NPROCESSORS_ONLN) linux; \
    mv -vfT ./out/minikube-linux-amd64 /export/minikube; \
    chmod +x /export/minikube;

FROM alpine:3.21 AS build-cri-dockerd
USER root
ENV CRI_DOCKERD_VERS=0.3.17
RUN \
    apk add --upgrade \
        bash \
        ca-certificates \
        msttcorefonts-installer \
        fontconfig \
        openssh-client \
        openssl \
        openssl-dev \
        alpine-sdk \
        autoconf \
        automake \
        bash \
        binutils-gold \
        build-base \
        cmake \
        file \
        file-dev \
        g++ \
        gcc \
        gnu-libiconv \
        gnu-libiconv-dev \
        gnupg \
        git \
        go \
        libgcc \
        libstdc++ \
        libtool \
        linux-headers \
        m4 \
        make \
        musl \
        musl-dev \
        nasm \
        pkgconf \
        pkgconf-dev \
        slibtool; \
    \
    # build cri-dockerd
    cd /home/; \
    git clone --depth=1 --branch=v$CRI_DOCKERD_VERS https://github.com/Mirantis/cri-dockerd.git ./cri-dockerd-$CRI_DOCKERD_VERS; \
    cd ./cri-dockerd-$CRI_DOCKERD_VERS; \
    make -j$(getconf _NPROCESSORS_ONLN) static-linux; \
    mkdir -pv /export/; \
    tar -xvzf ./packaging/static/build/linux/cri-dockerd-$CRI_DOCKERD_VERS.amd64.tgz -C ./ --strip 1; \
    mv -vfT ./cri-dockerd /export/cri-dockerd;

FROM scratch
COPY --from=build-kube /export/ /
COPY --from=build-cri-dockerd /export/ /
ENTRYPOINT ["/export/"]
