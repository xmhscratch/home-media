# final stage
FROM localhost:5000/sysenv:latest
WORKDIR /export/
ENV GO_ENV=development
ENV GIN_MODE=debug
ENV CC=/usr/bin/gcc
ENV CXX=/usr/bin/g++
ENV CC=/usr/bin/clang
ENV CXX=/usr/bin/clang++
ENV PKG_CONFIG=/usr/bin/pkg-config
ENV PKG_CONFIG_EXECUTABLE=/usr/bin/pkg-config
ENV LD_LIBRARY_PATH="/usr/lib:/usr/lib64:\$LD_LIBRARY_PATH"
ENV CGO_LDFLAGS="-L/usr/lib64"
ENV CGO_CFLAGS="-I/usr/include/ffmpeg"
ENV PKG_CONFIG_PATH="/usr/lib64/pkgconfig:\$PKG_CONFIG_PATH"
RUN \
    apk add --no-cache \
        --repository "http://dl-cdn.alpinelinux.org/alpine/edge/main" \
        --repository "http://dl-cdn.alpinelinux.org/alpine/edge/community" \
        wget2 \
        jq; \
    \
    apk add --virtual .build-deps \
        alpine-sdk \
    ; \
    ldconfig -v; \
    pkg-config --cflags -- libavcodec libavdevice libavfilter libavformat libswresample libswscale libavutil; \
    ffmpeg -version; \
    # cleanup
    apk del .build-deps; \
    rm -rf /var/cache/apk/*; \
    rm -rf /tmp/*;
EXPOSE 4000 4100 4110 4150 4200
ENTRYPOINT [ "/bin/sh", "-c" ]
# # CMD [ "nohup" "/export/hms > /var/log/hms.log 2>&1 &"]
