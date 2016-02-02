#!/bin/sh
# Make sure you install chromium build-deps first, e.g.
#
# $ sudo apt-get build-dep chromium-browser/sid
#

set -x
set -e

# Build (0) a single .a file, or (1, default) a bunch of *.so files.
# If 1, then you will need to set LD_LIBRARY_PATH at runtime or else install
# the *.so files as system libraries e.g. in /usr/lib/ or /usr/lib/${arch}
COMPONENT_BUILD="${COMPONENT_BUILD:-1}"
test "$COMPONENT_BUILD" = 1 -o "$COMPONENT_BUILD" = 0 || exit 2

################################################################################
# Download everything we need for the build.
# The result can be packed and shipped to offline reproducibility builders.

CHROMIUM_HOST=https://gsdview.appspot.com
DEBIAN_HOST=https://sources.debian.net
V_MAJOR="49"

# commented out for now; we don't need this when only building libjingle_peerconnection et al
#which gsutil || { echo >&2 'gsutil not found; try `pip install --user gsutil`;' exit 1; }

test -d webrtc/.git || git clone https://chromium.googlesource.com/external/webrtc

( cd webrtc
git config --replace-all remote.origin.fetch \
  +refs/branch-heads/*:refs/remotes/origin/branch-heads/* \
  '^\+refs/branch-heads/'
git fetch
git checkout branch-heads/"$V_MAJOR"
git reset --hard origin/branch-heads/"$V_MAJOR"
git clean -fdX

# download extra stuff listed in DEPS
test -d third_party/gflags/src || git clone --depth=1 \
  https://chromium.googlesource.com/external/gflags/src third_party/gflags/src
test -d third_party/junit-jar || git clone --depth=1 \
  https://chromium.googlesource.com/external/webrtc/deps/third_party/junit third_party/junit-jar

# download Google Cloud Storage resources
# commented out for now; we don't need this when only building libjingle_peerconnection et al
#find resources -name '*.sha1' -print0 | xargs -r0 -n1 \
#sh -c 'gsutil cp "gs://chromium-webrtc-resources/$(cat "$1")" "${1%.sha1}";
#if [ "$(sha1sum "${1%.sha1}" | cut "-d " -f1)" != "$(echo $(cat "$1"))" ]; then exit 255; fi' 0

# download Chromium release tarball
test -d chromium/src || ( cd chromium
CHROMIUM_PATH="$(wget -q -O- "$CHROMIUM_HOST/chromium-browser-official/?marker=chromium-$V_MAJOR" \
  | sed -nr -e 's,.*href="/(chromium-browser-official/chromium-'"$V_MAJOR"'\.[0-9\.]+-lite\.tar\.xz)".*,\1,gp' \
  | tail -n1)"
test -f "$(basename "$CHROMIUM_PATH")" || wget -nc "$CHROMIUM_HOST/$CHROMIUM_PATH"
CHROMIUM_DIR="$(tar xvf "$(basename "$CHROMIUM_PATH")" | sed -e 's@/.*@@' | uniq)"
ln -sf "$CHROMIUM_DIR" src )

# download Debian build file; use their config settings instead of figuring out our own
test -f chromium/rules || ( cd chromium
wget -N "$DEBIAN_HOST/$(wget -q -O- "$DEBIAN_HOST/src/chromium-browser/unstable/debian/rules/" | \
  sed -n -re 's/.*id="link_download" href="([^"]+)".*/\1/gp')" )

################################################################################
# Patches to upstream source code

# Patch chromium 49, upstream is just buggy
# sslv3 support dropped
# commented out for now; we don't need this when only building libjingle_peerconnection et al
#patch -Np1 <<EOF || true
#--- a/chromium/src/tools/swarming_client/third_party/requests/packages/urllib3/contrib/pyopenssl.py
#+++ b/chromium/src/tools/swarming_client/third_party/requests/packages/urllib3/contrib/pyopenssl.py
#@@ -39,7 +39,5 @@
# # Map from urllib3 to PyOpenSSL compatible parameter-values.
# _openssl_versions = {
#-    ssl.PROTOCOL_SSLv23: OpenSSL.SSL.SSLv23_METHOD,
#-    ssl.PROTOCOL_SSLv3: OpenSSL.SSL.SSLv3_METHOD,
#     ssl.PROTOCOL_TLSv1: OpenSSL.SSL.TLSv1_METHOD,
# }
# _openssl_verify = {
#EOF
# fix buggy upstream file that breaks the build
patch -Np1 <<EOF || true
--- a/chromium/src/chrome/test/data/webui_test_resources.grd
+++ b/chromium/src/chrome/test/data/webui_test_resources.grd
@@ -8,7 +8,6 @@
   </outputs>
   <release seq="1">
     <includes>
-      <include name="IDR_WEBUI_TEST_I18N_PROCESS_CSS_TEST" file="webui/i18n_process_css_test.html" flattenhtml="true" allowexternalscript="true" type="BINDATA" />
     </includes>
   </release>
 </grit>
EOF

# Patch debian 48 to work for upstream 49
grep -q use_sysroot=0 chromium/rules || sed -i -e '/^defines=/adefines+=use_sysroot=0' chromium/rules

################################################################################
# Add our own build targets

# Configure chromium
cat >chromium/src/Makefile <<'EOF'
include ../rules
.PHONY: configure_for_go_webrtc
configure_for_go_webrtc: override_dh_auto_configure
EOF

# Configure webrtc
cat >Makefile <<'EOF'
include chromium/rules
.PHONY: configure_for_go_webrtc
configure_for_go_webrtc:
{TAB_CHARACTER}GYP_DEFINES="$(defines) $(EXTRA_GYP_DEFINES)" python webrtc/build/gyp_webrtc
EOF
sed -i -e 's/{TAB_CHARACTER}/\t/g' Makefile

# Interim libjingle_peerconnection_internal target
patch -Np1 <<EOF
--- a/talk/libjingle.gyp
+++ b/talk/libjingle.gyp
@@ -43,6 +43,37 @@
     ['OS=="linux" or OS=="android"', {
       'targets': [
         {
+          # hacky target for go-webrtc. for this to work properly you MUST set
+          # CXXFLAGS='-fvisibility=default' when running gyp_webrtc
+          'target_name': 'libjingle_peerconnection_internal',
+          'type': 'shared_library',
+          'dependencies': [
+            '<(webrtc_root)/system_wrappers/system_wrappers.gyp:field_trial_default',
+            'libjingle_peerconnection',
+          ],
+          'sources': [
+            # need an empty file, otherwise gyp/ld goes crazy
+            'app/webrtc/peerconnection_internal.cc',
+          ],
+          'conditions': [
+            ['OS=="linux"', {
+              'defines': [
+                'HAVE_GTK',
+              ],
+              'conditions': [
+                ['use_gtk==1', {
+                  'link_settings': {
+                    'libraries': [
+                      '<!@(pkg-config --libs-only-l gobject-2.0 gthread-2.0'
+                          ' gtk+-2.0)',
+                    ],
+                  },
+                }],
+              ],
+            }],
+          ],
+        },
+        {
           'target_name': 'libjingle_peerconnection_jni',
           'type': 'static_library',
           'dependencies': [
EOF

patch -Np1 <<EOF
--- a/talk/app/webrtc/peerconnection_internal.cc
+++ b/talk/app/webrtc/peerconnection_internal.cc
@@ -0,0 +1 @@
+
EOF

) # exit webrtc dir

################################################################################
# The actual build

# build webrtc libs
( cd webrtc

touch ../.gclient # setup_links is overly strict, work around it
./setup_links.py -v

( cd chromium/src
rm -rf out
make configure_for_go_webrtc )

rm -rf out
[ $COMPONENT_BUILD = 1 ] && EXTRA_GYP_DEFINES="component=shared_library"
CXXFLAGS='-fvisibility=default' EXTRA_GYP_DEFINES="$EXTRA_GYP_DEFINES" \
  make configure_for_go_webrtc

ninja -C out/Release libjingle_peerconnection_internal
# archive all the headers. TODO: we can probably do less than this
tar cJf out/Release/libjingle_peerconnection_internal.headers.tar.xz $(find talk/ webrtc/ -type f -name '*.h')

if [ $COMPONENT_BUILD = 0 ]; then
# build a static library (non-thin) from libjingle_peerconnection_internal
( cd out/Release
rm -f lib/libjingle_peerconnection_internal.a
ninja -t query lib/libjingle_peerconnection_internal.so | \
  sed -n '/input: solink/,/outputs/p' | grep '\.a$' | \
  while read x; do ar rcs lib/libjingle_peerconnection_internal.a \
    $(ninja -t query "$x" | sed -n '/input: alink/,/outputs/p' | grep '\.o$'); done )
fi

) # exit webrtc dir

# copy webrtc libs to go-webrtc tree
rm -rf lib include
mkdir lib include
[ $COMPONENT_BUILD = 1 ] && cp webrtc/out/Release/lib/*.so lib || cp webrtc/out/Release/lib/*.a lib
tar -C include -xf webrtc/out/Release/libjingle_peerconnection_internal.headers.tar.xz

# build go-webrtc demo
go clean
if [ $COMPONENT_BUILD = 1 ]; then
    LD_LIBRARY_PATH="$PWD/lib"
    COMPONENT_LIBS="$(echo $(find lib -name *.so | sed -r -e 's,lib/lib(.*)\.so,-l\1,g'))"
else
    COMPONENT_LIBS="-ljingle_peerconnection_internal"
fi
sed -i -r -e 's/^(jingle_peerconnection_internal_libs)=.*/\1='"$COMPONENT_LIBS"'/g' \
  webrtc-linux-amd64.pc data/webrtc-data-linux-amd64.pc
LD_LIBRARY_PATH="$LD_LIBRARY_PATH" go build -v -x demo/chat/chat.go
