FROM localhost:5000/builder:latest AS build-apks
USER root
ENV GIT_SSL_NO_VERIFY=false
ENV ARCH=x86_64
COPY \
    manifest/apk-* \
    manifest/packages.txt \
    /tmp/build/
COPY dist/iso/.apks/ /export/iso/.apks/
RUN \
    cd /tmp/; \
    mkdir -pv /export/iso/; \
    mkdir -pv $HOME/tmp/; \
    ###########
    archdir="/export/iso/.apks/$ARCH"; \
    mkdir -pv "$archdir"; \
    apk update; \
    \
    pkgs=$(cat /tmp/build/packages.txt | tr '\n' ' '); \
    echo $pkgs; \
    apk fetch --link --recursive --output "$archdir" $pkgs; \
    if [[ $(find "$archdir"/*.apk -type f | wc -l) -gt 0 ]]; then \
        apk index \
            --allow-untrusted \
            --repository /export/iso/.apks/ \
            --rewrite-arch "$ARCH" \
            --index "$archdir"/APKINDEX.tar.gz \
            --output "$archdir"/APKINDEX.tar.gz \
            "$archdir"/*.apk; \
        abuild-sign "$archdir"/APKINDEX.tar.gz; \
    fi; \
    echo "/export/iso/.apks/" > /etc/apk/repositories; \
    apk update --allow-untrusted; \
    apk fix --upgrade --allow-untrusted alpine-keys;

FROM localhost:5000/builder:latest AS build-node-modules
USER root
ENV GIT_SSL_NO_VERIFY=false
ENV ARCH=x86_64
COPY package.json               /tmp/
COPY dist/iso/.node-modules/    /tmp/node_modules/
RUN \
    cd /tmp/; \
    mkdir -pv /export/iso/.node-modules/; \
    mkdir -pv $HOME/tmp/; \
    ###########
    npm install -g @angular/cli @nestjs/cli; \
    npm install --omit=dev; \
    mv -vf /tmp/node_modules/* /export/iso/.node-modules/;

FROM localhost:5000/builder:latest AS build-ci
USER root
ENV GIT_SSL_NO_VERIFY=false
ENV ARCH=x86_64
COPY \
    go.mod \
    package.json \
    angular.json \
    config.production.json \
    config.development.json \
    /tmp/build/ci/
COPY ci/ /tmp/build/ci/
RUN \
    apk add kustomize; \
    # ###########
    cd /tmp/; \
    mkdir -pv /export/iso/.ci/; \
    cat /tmp/build/ci/runtime.yml >> /tmp/build/ci/kustomization.yml; \
    kustomize build /tmp/build/ci/ --output /export/iso/.ci/hms-deploy.yml; \
    mv -vf \
        /tmp/build/ci/*-deploy.yml \
        /export/iso/.ci/; \
    \
    for fl in \
        hms \
        logstash \
        redis \
    ; do \
        svig_dir=/export/iso/.ci/$fl-svc-ingress.yml; \
        find /tmp/build/ci/$fl/ -type f -name service.yml | xargs -I @ sh -c "cat @ | tee -a $svig_dir; echo -e --- | tee -a $svig_dir"; \
        find /tmp/build/ci/$fl/ -type f -name ingress.yml | xargs -I @ sh -c "cat @ | tee -a $svig_dir; echo -e --- | tee -a $svig_dir"; \
    done;

FROM localhost:5000/builder:latest AS build-iso
USER root
ENV GIT_SSL_NO_VERIFY=false
ENV ARCH=x86_64
COPY \
    webos/*.sh \
    manifest/apk-* \
    manifest/packages.txt \
    manifest/kernel-modules.txt \
    /tmp/build/
COPY webos/script/                                  /tmp/build/script/
COPY channel/                                       /tmp/channel/
COPY dist/app/                                      /tmp/app/
COPY ci/*.sh                                        /tmp/bin/
COPY dist/bin/                                      /tmp/bin/
COPY dist/docker/preload-images.tar.gz              /tmp/docker/
COPY dist/docker/k3s-airgap-images-amd64.tar.zst    /tmp/docker/
COPY dist/iso/.node-modules/                        /tmp/node_modules/
COPY dist/iso/.apks/                                /tmp/apks/
COPY dist/iso/.ci/                                  /tmp/ci/
COPY dist/iso/node-modules.txt                      /tmp/
RUN \
    addgroup root abuild; \
    mkdir -pv $HOME/.mkimage/; \
    mkdir -pv /export/iso/; \
    mkdir -pv $HOME/tmp/; \
    abuild-keygen -i -a -n; \
    ###########
    cd $HOME/.mkimage/; \
    git clone --depth=1 --branch=3.21-stable https://gitlab.alpinelinux.org/alpine/aports.git ./aports; \
    sudo apk update; \
    export APORTS=$(realpath "$HOME/.mkimage/aports/"); \
    mv -vf \
        /tmp/build/genapkovl-hms.sh \
        /tmp/build/mkimg.hms.sh \
        $APORTS/scripts/; \
    \
    sudo chmod u+x \
        $APORTS/scripts/genapkovl-hms.sh \
        $APORTS/scripts/mkimg.hms.sh; \
    \
    ###########
    echo "/tmp/apks/" > /etc/apk/repositories; \
    $APORTS/scripts/mkimage.sh \
        --tag v3.21 \
        --outdir /export/iso/ \
        --arch $ARCH \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/main/ \
        --repository https://dl-cdn.alpinelinux.org/alpine/v3.21/community/ \
        --profile hms;

FROM scratch AS export-iso
COPY --from=build-iso /export/iso/ /
ENTRYPOINT ["/export/iso/"]

FROM scratch AS export-apks
COPY --from=build-apks /export/iso/ /
ENTRYPOINT ["/export/iso/"]

FROM scratch AS export-node-modules
COPY --from=build-node-modules /export/iso/ /
ENTRYPOINT ["/export/iso/"]

FROM scratch AS export-ci
COPY --from=build-ci /export/iso/ /
ENTRYPOINT ["/export/iso/"]
