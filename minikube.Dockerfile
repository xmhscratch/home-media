
FROM localhost:5000/gonode:latest AS build-kube
USER root
ENV KUBE_ARCH=x86_64
ENV KUBE_VERS=1.35.0
ENV GIT_SSL_NO_VERIFY=false
RUN apk update && apk add --no-cache \
        alpine-sdk \
        coreutils \
        wget \
        openssh \
        git \
        # ca-certificates-bundle \
        build-base \
        apk-tools \
        alpine-conf; \
    \
    # asdasd
    cd /home/; \
    git clone --depth=1 --branch=v$KUBE_VERS https://github.com/kubernetes/minikube.git ./minikube-$KUBE_VERS; \
    cd ./minikube-$KUBE_VERS; \
    make -j$(getconf _NPROCESSORS_ONLN) linux; \
    mv -vf ./out/minikube-linux-amd64 /export/minikube; \
    chmod +x /export/minikube;

FROM scratch
COPY --from=build-kube /export/ /
ENTRYPOINT ["/export/"]
