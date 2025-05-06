FROM localhost:5000/builder:latest AS build-kube
USER root
ENV KUBE_ARCH=x86_64
ENV KUBE_VERS=1.34.0
ENV GIT_SSL_NO_VERIFY=false
RUN \
    # build minikube
    cd /home/; \
    git clone --depth=1 --branch=v$KUBE_VERS https://github.com/kubernetes/minikube.git ./minikube-$KUBE_VERS; \
    cd ./minikube-$KUBE_VERS; \
    make -j$(getconf _NPROCESSORS_ONLN) linux; \
    mv -vfT ./out/minikube-linux-amd64 /export/minikube; \
    chmod +x /export/*;

FROM alpine:3.21 AS build-deps
USER root
ENV KUBERNETES_VERS=1.31.5
# ENV CRI_DOCKERD_VERS=0.3.17
ENV CNI_VERS=1.3.0
ENV CNI_PLUGIN_VERS=1.7.1
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
        curl \
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
    # # build cri-dockerd
    # cd /home/; \
    # git clone --depth=1 --branch=v$CRI_DOCKERD_VERS https://github.com/Mirantis/cri-dockerd.git ./cri-dockerd-$CRI_DOCKERD_VERS; \
    # cd ./cri-dockerd-$CRI_DOCKERD_VERS; \
    # make -j$(getconf _NPROCESSORS_ONLN) static-linux; \
    # mkdir -pv /export/; \
    # tar -xvzf ./packaging/static/build/linux/cri-dockerd-$CRI_DOCKERD_VERS.amd64.tgz -C ./ --strip 1; \
    # mv -vfT ./cri-dockerd /export/cri-dockerd; \
    # build cni plugins
    cd /home/; \
    git clone --depth=1 --branch=v$CNI_PLUGIN_VERS https://github.com/containernetworking/plugins.git ./plugins-$CNI_PLUGIN_VERS; \
    cd ./plugins-$CNI_PLUGIN_VERS; \
    go get github.com/containernetworking/cni; \
    ./build_linux.sh; \
    mkdir -pv /export/; \
    mv -vf ./bin/* /export/; \
    # download kube executable
    for exe in \
        kubeadm \
        kubectl \
        kubelet \
    ; do \
        curl -sS --skip-existing https://cdn.dl.k8s.io/release/v$KUBERNETES_VERS/bin/linux/amd64/"$exe" -o /export/"$exe"; \
    done; \
    \
    chmod +x /export/*;

FROM scratch
COPY --from=build-kube /export/ /
COPY --from=build-deps /export/ /
ENTRYPOINT ["/export/"]
