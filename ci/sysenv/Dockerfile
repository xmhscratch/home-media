FROM localhost:5000/gonode:latest
WORKDIR /tmp
# ENV FFMPEG_VERSION=n5.1.6
ENV FFMPEG_VERSION=n7.1
ENV NV_CODEC_HEADERS_VERSION=n12.2.72.0
RUN \
    apk update && apk add \
        bash \
        ca-certificates \
        # curl \
        # doas \
        # fontconfig \
        # linux-headers \
        # m4 \
        # musl \
        # musl-dbg \
        # musl-utils \
        # openssh-client \
        python3 \
        py3-pip \
        tar \
        # zip \
        unzip \
        wget \
    ; \
    apk add \
    #     chrpath \
    #     ghostscript \
    #     gtk+2.0 \
    #     inotify-tools \
    #     jpeg \
    #     lcms2 \
    #     leptonica \
        libass \
    #     libexif \
    #     libffi \
    #     libgphoto2 \
    #     libjpeg \
    #     libjpeg-turbo \
    #     libltdl \
    #     libpcrecpp \
    #     libpng \
    #     libraw1394 \
    #     librsvg \
    #     libunwind \
        libvorbis \
        libwebp \
    #     libwmf \
    #     libx11 \
    #     libxext \
    #     libxml2 \
    #     libxt \
    #     openblas \
    #     openexr-libopenexrcore \
    #     pcre \
    #     pcre-tools \
    #     potrace \
    #     tesseract-ocr \
    #     shaderc \
        libsrt \
    #     openal-soft \
    #     opencl \
    #     ffnvcodec-headers \
    #     vulkan-headers \
    #     lilv \
    #     zlib \
    ; \
    apk add --virtual .libs-main-dev \
    #     bzip2-dev \
    #     clang-dev \
    #     dpkg-dev \
    #     fftw-dev \
    #     fontconfig-dev \
    #     freetype-dev \
    #     ghostscript-dev \
    #     jpeg-dev \
    #     lcms2-dev \
    #     leptonica-dev \
        libass-dev \
    #     libexif-dev \
    #     libffi-dev \
    #     libgphoto2-dev \
    #     libheif-dev \
    #     libjpeg-turbo-dev \
    #     libpng-dev \
    #     libraw1394-dev \
    #     librsvg-dev \
        libvorbis-dev \
        libwebp-dev \
    #     libwmf-dev \
    #     libx11-dev \
    #     libxext-dev \
    #     libxml2-dev \
    #     libxt-dev \
    #     musl-dev \
    #     openblas-dev \
    #     openjpeg-dev \
    #     pango-dev \
    #     pcre-dev \
    #     perl-dev \
    #     potrace-dev \
    #     python3-dev \
    #     tesseract-ocr-dev \
    #     tiff-dev \
    #     shaderc-dev \
        libsrt-dev \
    #     openal-soft-dev \
    #     opencl-dev \
    #     lilv-dev \
    #     zlib-dev \
    ; \
    apk add --repository "http://dl-cdn.alpinelinux.org/alpine/edge/main" --repository "http://dl-cdn.alpinelinux.org/alpine/edge/community" \
        # dav1d \
        # frei0r-plugins \
        # git \
        # glog \
        # gflags \
        # gnutls \
        # ladspa \
        # libbluray \
        # libdrm \
        # libmodplug \
        # opencore-amr \
        # libopenmpt \
        opus \
        # rubberband \
        # libssh2 \
        # speex \
        # libtheora \
        # vidstab \
        libvpx \
        x264 \
        x265 \
        # xvidcore \
        # zimg \
        # libcdio \
        # libcdio-paranoia \
        # librsvg \
        ffmpeg-libavcodec \
        ffmpeg-libavdevice \
        ffmpeg-libavfilter \
        ffmpeg-libavformat \
        ffmpeg-libswresample \
        ffmpeg-libswscale \
        ffmpeg-libavutil \
        # samurai \
    ; \
    apk add --virtual .libs-edge-dev --repository "http://dl-cdn.alpinelinux.org/alpine/edge/main" --repository "http://dl-cdn.alpinelinux.org/alpine/edge/community" \
    #     dav1d-dev \
    #     frei0r-plugins-dev \
    #     glog-dev \
    #     gflags-dev \
    #     gnutls-dev \
    #     ladspa-dev \
    #     libbluray-dev \
    #     libdrm-dev \
    #     libmodplug-dev \
    #     opencore-amr-dev \
    #     libopenmpt-dev \
        opus-dev \
    #     rubberband-dev \
    #     libssh2-dev \
    #     speex-dev \
    #     libtheora-dev \
    #     vidstab-dev \
        libvpx-dev \
        x264-dev \
        x265-dev \
    #     xvidcore-dev \
    #     zimg-dev \
    #     libcdio-dev \
    #     libcdio-paranoia-dev \
    #     librsvg-dev \
    ; \
    apk add --virtual .build-dev \
        file-dev \
        openssl-dev \
        valgrind-dev \
        gnu-libiconv-dev \
        glib-dev \
        pkgconf-dev \
    ; \
    apk add --virtual .build \
        file \
        nasm \
        ninja \
        git \
        git-lfs \
        slibtool \
        valgrind \
        autoconf \
        automake \
        alpine-sdk \
        binutils-gold \
        build-base \
        cmake \
        gnu-libiconv \
        gnupg \
        glib \
        libgcc \
        libtool \
        libstdc++ \
        make \
        pkgconf \
    ; \
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
    # \
    export PYTHON3_EXECUTABLE=$(which python3); \
    export PYTHON3_INCLUDE_DIR=$(python3 -c "from distutils.sysconfig import get_python_inc; print(get_python_inc())"); \
    export PYTHON3_INCLUDE_PATH=$PYTHON3_INCLUDE_DIR; \
    export PYTHON3_VERSION=$(echo $PYTHON3_INCLUDE_PATH | rev | cut -d '/' -f 1 | rev); \
    export PYTHON3_LIBRARY=$(python3 -c "import distutils.sysconfig as sysconfig; print(sysconfig.get_config_var('LIBDIR'))"); \
    export PYTHON3_LIBRARIES=$PYTHON3_LIBRARY; \
    export PYTHON3_NUMPY_VERSION=$(python3 -c "import numpy; print(numpy.__version__)"); \
    export PYTHON3_NUMPY_INCLUDE_DIRS=$PYTHON3_LIBRARY/$PYTHON3_VERSION/site-packages/numpy/core/include; \
    export PYTHONPATH=$PYTHON3_LIBRARY/$PYTHON3_VERSION/site-packages:$PYTHONPATH; \
    # \
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
        # --enable-libopencore-amrnb \
        # --enable-libopencore-amrwb \
        --disable-libvo-amrwbenc \
        # --enable-version3 \
        # --enable-bzlib \
        # --enable-fontconfig \
        # --enable-frei0r \
        --disable-gcrypt \
        # --enable-gnutls \
        # --enable-ladspa \
        # --enable-libaom \
        # --enable-libdav1d \
        --enable-libass \
        # --enable-libbluray \
        # --enable-libcdio \
        # --enable-libdrm \
        # --enable-libjack \
        --disable-libjxl \
        # --enable-libfreetype \
        # --enable-libfribidi \
        --disable-libgsm \
        --disable-libmp3lame \
        --disable-libmysofa \
        --enable-nvenc \
        # --enable-openal \
        # --enable-opencl \
        --disable-opengl \
        # --enable-libopenjpeg \
        # --enable-libopenmpt \
        --enable-libopus \
        --disable-libpulse \
        # --enable-librsvg \
        # --enable-librubberband \
        --disable-libsoxr \
        # --enable-libspeex \
        --enable-libsrt \
        --disable-libssh \
        --disable-libsvtav1 \
        # --enable-libtheora \
        --enable-libvorbis \
        --disable-libv4l2 \
        # --enable-libvidstab \
        --disable-libvmaf \
        # --enable-version3 \
        --disable-vapoursynth \
        --enable-libvpx \
        --disable-vulkan \
        # --enable-libshaderc \
        --enable-libx264 \
        --enable-libx265 \
        # --enable-libxvid \
        # --enable-libxml2 \
        # --enable-libzimg \
        --disable-libzvbi \
        # --enable-lv2 \
        # --enable-avfilter \
        # --enable-libmodplug \
        # --enable-postproc \
        # --enable-pthreads \
        --enable-static \
        # --enable-shared \
        --enable-gpl \
        --disable-debug \
        --disable-stripping \
        --shlibdir=/usr/lib64 \
        # --enable-lto \
        --disable-libmfx \
        # --enable-runtime-cpudetect \
        --disable-doc \
        --disable-htmlpages \
        --disable-manpages \
        --disable-podpages \
        --disable-txtpages \
    ; \
    make -j$(getconf _NPROCESSORS_ONLN); \
    make install; \
    ldconfig; \
    pkg-config --cflags -- libavcodec libavdevice libavfilter libavformat libswresample libswscale libavutil; \
    ffmpeg -version; \
    \
    # # install imagemagick
    # apk add \
    #     imagemagick \
    #     imagemagick-libs \
    #     imagemagick-c++ \
    # ; \
    # ldconfig; \
    # pkg-config --cflags --libs MagickWand; \
    # \
    # cleanup
    deactivate; \
    apk del \
        .build \
        .build-dev \
        .libs-main-dev \
        .libs-edge-dev \
    ; \
    rm -rf /var/cache/apk/*; \
    rm -rf /tmp/*;

WORKDIR /export/
ENTRYPOINT [ "/bin/sh", "-c" ]
