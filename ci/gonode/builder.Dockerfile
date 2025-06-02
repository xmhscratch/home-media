FROM localhost:5000/gonode:latest
USER root
ENV GIT_SSL_NO_VERIFY=false
RUN apk update && apk add --no-cache \
        alpine-sdk \
        build-base \
        apk-tools \
        alpine-conf \
        ca-certificates \
        busybox \
        fakeroot \
        xorriso \
        squashfs-tools \
        sudo \
        mtools \
        dosfstools \
        grub-efi \
        python3 \
        py3-pylspci; \
    \
    # cleanup
    rm -rf /var/cache/apk/*; \
    rm -rf /tmp/*; \
    apk update;
