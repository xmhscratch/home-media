# # build stage
# FROM localhost:5000/sysenv:latest AS builder

# WORKDIR /export/

# ENV GOOS linux
# ENV GOARCH amd64
# ENV CGO_ENABLED 0 

# # build
# go get -d -v ./...; \
# go install -v ./...; \
# go mod vendor; \
# go build -ldflags="-s -w" -mod=vendor -o /go/bin/api ./cmd/api/main.go; \
# go build -ldflags="-s -w" -mod=vendor -o /go/bin/downloader ./cmd/downloader/main.go; \
# go build -ldflags="-s -w" -mod=vendor -o /go/bin/encoder ./cmd/encoder/main.go; \
# go build -ldflags="-s -w" -mod=vendor -o /go/bin/file ./cmd/file/main.go; \
# # build
# npm install; \
# npm run client:build; \
# npm run backend:build; \

# final stage
FROM localhost:5000/sysenv:latest

WORKDIR /export/

ENV GO_ENV development
ENV GIN_MODE debug
ENV CC=/usr/bin/gcc
ENV CXX=/usr/bin/g++
ENV CC=/usr/bin/clang
ENV CXX=/usr/bin/clang++
ENV PKG_CONFIG=/usr/bin/pkg-config
ENV PKG_CONFIG_EXECUTABLE=/usr/bin/pkg-config
ENV LD_LIBRARY_PATH=/usr/lib:/usr/lib64:$LD_LIBRARY_PATH
ENV CGO_LDFLAGS="-L/usr/lib64"
ENV CGO_CFLAGS="-I/usr/include/ffmpeg"
ENV PKG_CONFIG_PATH="/usr/lib64/pkgconfig:$PKG_CONFIG_PATH"

# COPY --from=builder \
#     /go/bin/api \
#     /go/bin/downloader \
#     /go/bin/encoder \
#     /go/bin/file \
#     /export/cmd/
# COPY --from=builder \
#     /export/dist/ \
#     /export/dist/
COPY \
    .env \
    .browserslistrc \
    .env \
    angular.json \
    config.development.json \
    config.production.json \
    go.mod \
    go.sum \
    package.json \
    postcss.config.js \
    tailwind.config.js \
    tsconfig.app.json \
    tsconfig.json \
    tsconfig.spec.json \
    /export/
COPY ci/*.sh /export/bin/

RUN \
    apk add --upgrade --no-cache --allow-untrusted --update-cache \
        --repository "http://dl-cdn.alpinelinux.org/alpine/edge/main" \
        --repository "http://dl-cdn.alpinelinux.org/alpine/edge/community" \
        wget2 \
        jq; \
    \
    apk add --virtual .build-deps \
        file \
        file-dev \
        git-lfs \
        nano \
        nasm \
        ninja \
        openssh-client \
        openssl-dev \
        slibtool \
        valgrind \
        valgrind-dev \
        autoconf \
        automake \
        alpine-sdk \
        binutils-gold \
        build-base \
        cmake \
        gnu-libiconv \
        gnu-libiconv-dev \
        gnupg \
        glib \
        glib-dev \
        libgcc \
        libtool \
        linux-headers \
        libstdc++ \
        make \
        pkgconf \
        pkgconf-dev; \
    \
    ldconfig -v; \
    pkg-config --cflags -- libavcodec libavdevice libavfilter libavformat libswresample libswscale libavutil; \
    ffmpeg -version; \
    # permissions
    # chmod u+x /export/cmd/*; \
    chmod u+x /export/bin/*.sh; \
    # cleanup
    apk del .build-deps; \
    rm -rf /var/cache/apk/*; \
    rm -rf /tmp/*;

EXPOSE 4000 4100 4110 4150 4200
ENTRYPOINT [ "/bin/sh", "-c" ]
# CMD [ "nohup" "/export/hms > /var/log/hms.log 2>&1 &"]
