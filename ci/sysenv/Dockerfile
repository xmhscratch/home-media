FROM localhost:5000/gonode:latest

WORKDIR /tmp
# ENV FFMPEG_VERSION=n5.1.6
ENV FFMPEG_VERSION=n7.1
ENV NV_CODEC_HEADERS_VERSION=n12.2.72.0

RUN \
    apk -U upgrade; \
    \
    apk add --upgrade --no-cache \
        bash \
        bzip2-dev \
        ca-certificates \
        chrpath \
        clang-dev \
        curl \
        doas \
        dpkg-dev \
        fftw-dev \
        fontconfig \
        fontconfig-dev \
        freetype-dev \
        ghostscript \
        ghostscript-dev \
        git \
        git-lfs \
        gtk+2.0 \
        inotify-tools \
        jpeg \
        jpeg-dev \
        lcms2 \
        lcms2-dev \
        leptonica \
        leptonica-dev \
        libass \
        libass-dev \
        libexif \
        libexif-dev \
        libffi \
        libffi-dev \
        libgphoto2 \
        libgphoto2-dev \
        libheif-dev \
        libjpeg \
        libjpeg-turbo \
        libjpeg-turbo-dev \
        libltdl \
        libpcrecpp \
        libpng \
        libpng-dev \
        libraw1394 \
        libraw1394-dev \
        librsvg \
        librsvg-dev \
        libunwind \
        libwebp \
        libwebp-dev \
        libwmf \
        libwmf-dev \
        libx11 \
        libx11-dev \
        libxext \
        libxext-dev \
        libxml2 \
        libxml2-dev \
        libxt \
        libxt-dev \
        linux-headers \
        m4 \
        msttcorefonts-installer \
        musl \
        musl-dbg \
        musl-dev \
        musl-utils \
        openblas \
        openblas-dev \
        openjpeg-dev \
        openssh-client \
        openssl \
        openexr-libopenexrcore \
        pango-dev \
        pcre \
        pcre-dev \
        pcre-tools \
        perl-dev \
        potrace \
        potrace-dev \
        python3 \
        python3-dev \
        py3-pip \
        tar \
        tesseract-ocr \
        tesseract-ocr-dev \
        tiff-dev \
        unzip \
        wget \
        shaderc \
        shaderc-dev \
        libsrt \
        libsrt-dev \
        openal-soft \
        openal-soft-dev \
        opencl \
        opencl-dev \
        ffnvcodec-headers \
        vulkan-headers \
        lilv \
        lilv-dev \
        zip \
        zlib-dev \
        zlib; \
    \
    apk add --upgrade --no-cache \
        ccache@3.14 \
        hdf5-dev@3.14 \
        hdf5@3.14 \
        libexecinfo@3.16 \
        libexecinfo-dev@3.16 \
        libjxl@3.16 \
        libjxl-dev@3.16 \
        libjxl-pixbuf-loader@3.16 \
        jasper@3.10 \
        jasper-dev@3.10; \
    \
    apk add --upgrade --no-cache --allow-untrusted --update-cache \
            --repository "http://dl-cdn.alpinelinux.org/alpine/edge/main" \
            --repository "http://dl-cdn.alpinelinux.org/alpine/edge/community" \
        dav1d \
        dav1d-dev \
        frei0r-plugins \
        frei0r-plugins-dev \
        git \
        glog \
        glog-dev \
        gflags \
        gflags-dev \
        gnutls \
        gnutls-dev \
        ladspa \
        ladspa-dev \
        libbluray \
        libbluray-dev \
        libdrm \
        libdrm-dev \
        libmodplug \
        libmodplug-dev \
        opencore-amr \
        opencore-amr-dev \
        libopenmpt \
        libopenmpt-dev \
        opus \
        opus-dev \
        rubberband \
        rubberband-dev \
        libssh2 \
        libssh2-dev \
        speex \
        speex-dev \
        libtheora \
        libtheora-dev \
        vidstab \
        vidstab-dev \
        libvpx \
        libvpx-dev \
        x264 \
        x264-dev \
        xvidcore \
        xvidcore-dev \
        zimg \
        zimg-dev \
        libcdio \
        libcdio-dev \
        libcdio-paranoia \
        libcdio-paranoia-dev \
        librsvg \
        librsvg-dev \
        ffmpeg-libavcodec \
        ffmpeg-libavdevice \
        ffmpeg-libavfilter \
        ffmpeg-libavformat \
        ffmpeg-libswresample \
        ffmpeg-libswscale \
        ffmpeg-libavutil \
        samurai; \
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
    # setup environment
    export CC=/usr/bin/gcc; \
    export CXX=/usr/bin/g++; \
    export CC=/usr/bin/clang; \
    export CXX=/usr/bin/clang++; \
    export PKG_CONFIG=/usr/bin/pkg-config; \
    export PKG_CONFIG_PATH=$(pkg-config --variable pc_path pkg-config):$PKG_CONFIG_PATH; \
    export PKG_CONFIG_EXECUTABLE=/usr/bin/pkg-config; \
    export GLIB_CFLAGS=$(pkg-config --cflags glib-2.0); \
    export GLIB_LIBS=$(pkg-config --libs glib-2.0); \
    export LD_LIBRARY_PATH=/usr/lib:/usr/lib64:$LD_LIBRARY_PATH; \
    ln -s /usr/include/locale.h /usr/include/xlocale.h; \
    ldconfig -v; \
    # python3
    python3 -m venv /root/.venv; \
    . /root/.venv/bin/activate; \
    python3 -m pip install setuptools wheel; \
    python3 -m pip install --force-reinstall --upgrade numpy; \
    deactivate; \
    \
    export PYTHON3_EXECUTABLE=$(which python3); \
    export PYTHON3_INCLUDE_DIR=$(python3 -c "from distutils.sysconfig import get_python_inc; print(get_python_inc())"); \
    export PYTHON3_INCLUDE_PATH=$PYTHON3_INCLUDE_DIR; \
    export PYTHON3_VERSION=$(echo $PYTHON3_INCLUDE_PATH | rev | cut -d '/' -f 1 | rev); \
    export PYTHON3_LIBRARY=$(python3 -c "import distutils.sysconfig as sysconfig; print(sysconfig.get_config_var('LIBDIR'))"); \
    export PYTHON3_LIBRARIES=$PYTHON3_LIBRARY; \
    export PYTHON3_NUMPY_VERSION=$(python3 -c "import numpy; print(numpy.__version__)"); \
    export PYTHON3_NUMPY_INCLUDE_DIRS=$PYTHON3_LIBRARY/$PYTHON3_VERSION/site-packages/numpy/core/include; \
    export PYTHONPATH=$PYTHON3_LIBRARY/$PYTHON3_VERSION/site-packages:$PYTHONPATH; \
    \
    # install ffmpeg nv-codec-headers
    cd /tmp/; \
    wget --progress=bar:force:noscroll https://github.com/FFmpeg/nv-codec-headers/archive/refs/tags/n12.2.72.0.zip -O /tmp/nv-codec-headers-$NV_CODEC_HEADERS_VERSION.zip; \
    unzip -q nv-codec-headers-$NV_CODEC_HEADERS_VERSION.zip; \
    cd nv-codec-headers-$NV_CODEC_HEADERS_VERSION; \
    make -j$(getconf _NPROCESSORS_ONLN); \
    make install; \
    # install ffmpeg
    cd /tmp/; \
    wget --progress=bar:force:noscroll https://github.com/FFmpeg/FFmpeg/archive/refs/tags/$FFMPEG_VERSION.zip -O /tmp/FFmpeg-$FFMPEG_VERSION.zip; \
    unzip -q FFmpeg-$FFMPEG_VERSION.zip; \
    cd ./FFmpeg-$FFMPEG_VERSION; \
    export CGO_LDFLAGS="-L/usr/lib64"; \
    export CGO_CFLAGS="-I/usr/include/ffmpeg"; \
    export PKG_CONFIG_PATH="/usr/lib64/pkgconfig:$PKG_CONFIG_PATH"; \
    ./configure \
        --prefix=/usr \
        --bindir=/usr/bin \
        --datadir=/usr/share/ffmpeg \
        --docdir=/usr/share/doc/ffmpeg \
        --incdir=/usr/include/ffmpeg \
        --libdir=/usr/lib64 \
        --mandir=/usr/share/man \
        --arch=x86_64 \
        --optflags='-O2 -flto=auto -ffat-lto-objects -fexceptions -g -grecord-gcc-switches -pipe -Wall -Werror=format-security -Wp,-D_FORTIFY_SOURCE=2 -Wp,-D_GLIBCXX_ASSERTIONS -specs=/usr/lib/rpm/redhat/redhat-hardened-cc1 -fstack-protector-strong -specs=/usr/lib/rpm/redhat/redhat-annobin-cc1 -m64 -march=x86-64-v2 -mtune=generic -fasynchronous-unwind-tables -fstack-clash-protection -fcf-protection' \
        --enable-libopencore-amrnb \
        --enable-libopencore-amrwb \
        --disable-libvo-amrwbenc \
        --enable-version3 \
        --enable-bzlib \
        --enable-fontconfig \
        --enable-frei0r \
        --disable-gcrypt \
        --enable-gnutls \
        --enable-ladspa \
        --enable-libaom \
        --enable-libdav1d \
        --enable-libass \
        --enable-libbluray \
        --enable-libcdio \
        --enable-libdrm \
        --enable-libjack \
        --disable-libjxl \
        --enable-libfreetype \
        --enable-libfribidi \
        --disable-libgsm \
        --disable-libmp3lame \
        --disable-libmysofa \
        --enable-nvenc \
        --enable-openal \
        --enable-opencl \
        --disable-opengl \
        --enable-libopenjpeg \
        --enable-libopenmpt \
        --enable-libopus \
        --disable-libpulse \
        --enable-librsvg \
        --enable-librubberband \
        --disable-libsoxr \
        --enable-libspeex \
        --enable-libsrt \
        --disable-libssh \
        --disable-libsvtav1 \
        --enable-libtheora \
        --enable-libvorbis \
        --disable-libv4l2 \
        --enable-libvidstab \
        --disable-libvmaf \
        --enable-version3 \
        --disable-vapoursynth \
        --enable-libvpx \
        --disable-vulkan \
        --enable-libshaderc \
        --enable-libx264 \
        --enable-libx265 \
        --enable-libxvid \
        --enable-libxml2 \
        --enable-libzimg \
        --disable-libzvbi \
        --enable-lv2 \
        --enable-avfilter \
        --enable-libmodplug \
        --enable-postproc \
        --enable-pthreads \
        --enable-static \
        --enable-shared \
        --enable-gpl \
        --disable-debug \
        --disable-stripping \
        --shlibdir=/usr/lib64 \
        --enable-lto \
        --disable-libmfx \
        --enable-runtime-cpudetect; \
    make -j$(getconf _NPROCESSORS_ONLN); \
    make install; \
    ldconfig; \
    pkg-config --cflags -- libavcodec libavdevice libavfilter libavformat libswresample libswscale libavutil; \
    ffmpeg -version; \
    # install imagemagick
    apk add --upgrade \
        imagemagick \
        imagemagick-dev \
        imagemagick-libs \
        imagemagick-c++; \
    \
    ldconfig; \
    pkg-config --cflags --libs MagickWand; \
    \
    # cleanup
    apk del .build-deps; \
    rm -rf /var/cache/apk/*; \
    rm -rf /tmp/*;

RUN \
    go env -w GOPRIVATE=localhost; \
    git config --global user.name "root"; \
    git config --global user.email "root@com.home-media.me"; \
    git config --global http.sslverify "false";

WORKDIR /export/
ENTRYPOINT [ "/bin/sh", "-c" ]
