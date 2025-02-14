# # build stage
# FROM localhost:5000/sysenv:latest AS builder
# WORKDIR /go/src/home-media/
# # ENV GO111MODULE on
# # ENV GIT_SSL_NO_VERIFY 1
# COPY . .
# RUN \
#     apk add --virtual .build-deps \
#         file \
#         file-dev \
#         git-lfs \
#         nano \
#         nasm \
#         ninja \
#         openssh-client \
#         openssl-dev \
#         slibtool \
#         valgrind \
#         valgrind-dev \
#         autoconf \
#         automake \
#         alpine-sdk \
#         binutils-gold \
#         build-base \
#         cmake \
#         gnu-libiconv \
#         gnu-libiconv-dev \
#         gnupg \
#         glib \
#         glib-dev \
#         libgcc \
#         libtool \
#         linux-headers \
#         libstdc++ \
#         make \
#         pkgconf \
#         pkgconf-dev; \
#     \
#     export CC=/usr/bin/gcc; \
#     export CXX=/usr/bin/g++; \
#     export CC=/usr/bin/clang; \
#     export CXX=/usr/bin/clang++; \
#     export PKG_CONFIG=/usr/bin/pkg-config; \
#     export PKG_CONFIG_EXECUTABLE=/usr/bin/pkg-config; \
#     export PKG_CONFIG_PATH=$(pkg-config --variable pc_path pkg-config):$PKG_CONFIG_PATH; \
#     export GLIB_CFLAGS=$(pkg-config --cflags glib-2.0); \
#     export GLIB_LIBS=$(pkg-config --libs glib-2.0); \
#     export LD_LIBRARY_PATH=/usr/lib:/usr/lib64:$LD_LIBRARY_PATH; \
#     export CGO_LDFLAGS="-L/usr/lib64"; \
#     export CGO_CFLAGS="-I/usr/include/ffmpeg"; \
#     export PKG_CONFIG_PATH="/usr/lib64/pkgconfig:$PKG_CONFIG_PATH"; \
#     ln -s /usr/include/locale.h /usr/include/xlocale.h; \
#     ldconfig -v; \
#     pkg-config --cflags -- libavcodec libavdevice libavfilter libavformat libswresample libswscale libavutil; \
#     ffmpeg -version; \
#     go get -d -v ./...; \
#     go install -v ./...; \
#     go mod vendor; \
#     go build -mod=vendor -o /go/bin/hms; \
#     # cleanup
#     apk del .build-deps; \
#     rm -rf /var/cache/apk/*; \
#     rm -rf /tmp/*;

# final stage
FROM localhost:5000/sysenv:latest

# WORKDIR /export/
WORKDIR /go/src/home-media/

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

# COPY --from=builder /go/bin/hms /export/hms
COPY ["*.sh", "/bin/home-media/"]

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
    chmod u+x /bin/home-media/*.sh;

EXPOSE 4110 4150
ENTRYPOINT [ "/bin/sh", "-c" ]
# CMD [ "nohup" "/export/hms > /var/log/hms.log 2>&1 &"]
