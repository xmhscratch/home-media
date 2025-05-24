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

FROM localhost:5000/builder:latest AS build-iso
USER root
ENV GIT_SSL_NO_VERIFY=false
ENV ARCH=x86_64
COPY \
    webos/*.sh \
    manifest/apk-* \
    manifest/packages.txt \
    /tmp/build/
        
COPY package.json                                   /tmp/
COPY ci/                                            /tmp/build/ci/
COPY webos/script/                                  /tmp/build/script/
COPY channel/                                       /tmp/channel/
COPY dist/app/                                      /tmp/app/
COPY dist/bin/                                      /tmp/bin/
COPY dist/docker/preload-images.tar.gz              /tmp/docker/
COPY dist/docker/k3s-airgap-images-amd64.tar.zst    /tmp/docker/
COPY dist/iso/.node-modules/                        /tmp/node_modules/
COPY dist/iso/.apks/                                /export/iso/.apks/
RUN \
    apk add kustomize; \
    ###########
    addgroup root abuild; \
    mkdir -pv $HOME/.mkimage/; \
    mkdir -pv /export/iso/; \
    mkdir -pv $HOME/tmp/; \
    abuild-keygen -i -a -n; \
    ###########
    cd /tmp/; \
    mkdir -pv /tmp/ci/; \
    mv -vf \
        /tmp/build/ci/*-deploy.yml \
        /tmp/ci/; \
    \
    kustomize build /tmp/build/ci/ --output /tmp/ci/hms-deploy.yml; \
    for fl in \
        logstash \
        redis \
        hms \
    ; do \
        svig_dir=/tmp/ci/$fl-svc-ingress.yml; \
        find /tmp/build/ci/$fl/service.yml -type f | xargs -I @ sh -c "cat @ | tee -a $svig_dir; echo -e --- | tee -a $svig_dir" >/dev/null &2>1; \
        find /tmp/build/ci/$fl/ingress.yml -type f | xargs -I @ sh -c "cat @ | tee -a $svig_dir; echo -e --- | tee -a $svig_dir" >/dev/null &2>1; \
    done; \
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
