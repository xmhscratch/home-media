DESCRIPTION = "Ninja is a small build system with a focus on speed."
LICENSE = "Apache-2"

inherit native

# LIC_FILES_CHKSUM = "file://COPYING;md5=a81586a64ad4e476c791cda7e2f2c52e"
# SRCREV="2daa09ba270b0a43e1929d29b073348aa985dfaa"
# SRCBRANCH="release"
# SRC_URI = "git://github.com/martine/ninja.git;protocol=https;branch=${SRCBRANCH}"

LIC_FILES_CHKSUM = "file://COPYING;md5=a81586a64ad4e476c791cda7e2f2c52e"
SRC_URI = "file:///home/web/repos/qtwebengine-chromium"

S="${WORKDIR}/ninja"

do_compile() {
    python ${S}/bootstrap.py
}

do_install() {
    install -d ${D}${bindir}
    install -m 0755 ${S} ${D}${bindir}/ninja
}
