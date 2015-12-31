#!/usr/bin/env bash

PROJECT_DIR=$(pwd)
THIRD_PARTY_DIR="$PROJECT_DIR/third_party"
WEBRTC_REPO="https://chromium.googlesource.com/external/webrtc"
WEBRTC_DIR="$THIRD_PARTY_DIR/webrtc"
WEBRTC_SRC="$WEBRTC_DIR/src"
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
CONFIG="Debug"
COMMIT="bd7d8f7e2b824a887aa12236cb6185d446d7da61"

INCLUDE_DIR="$PROJECT_DIR/include"
LIB_DIR="$PROJECT_DIR/lib"

# TODO(arlolra): depot_tools

GYP_DEFINES="include_tests=0"

mkdir -p $THIRD_PARTY_DIR
mkdir -p $INCLUDE_DIR
mkdir -p $LIB_DIR

if [[ -d $WEBRTC_DIR ]]; then
	echo "Sync'ing webrtc ..."
	pushd $WEBRTC_DIR
	gclient sync
	popd
else
	echo "Getting webrtc ..."
	mkdir -p $WEBRTC_DIR
	pushd $WEBRTC_DIR
	gclient config --name src $WEBRTC_REPO
	gclient sync
	popd
fi

echo "Checking out latest tested / compatible version of webrtc ..."
pushd $WEBRTC_SRC
git checkout $COMMIT
popd

echo "Generating build scripts ..."
pushd $WEBRTC_SRC
python webrtc/build/gyp_webrtc
popd

echo "Building webrtc ..."
pushd $WEBRTC_SRC
ninja -C out/$CONFIG
popd

echo "Copying headers ..."
pushd $WEBRTC_SRC
for h in $(find talk/ webrtc/ -type f -name '*.h')
do
	mkdir -p "$INCLUDE_DIR/$(dirname $h)"
	cp $h "$INCLUDE_DIR/$h"
done
popd

echo "Concatenating libraries ..."
pushd $WEBRTC_SRC/out/$CONFIG
ls *.a > filelist
libtool -static -o libwebrtc-magic.a -filelist filelist
# on osx:
# strip -S -x -o libwebrtc-magic.a libwebrtc-magic.a
mv libwebrtc-magic.a $LIB_DIR/libwebrtc-$OS-$ARCH-magic.a
popd

echo "Build complete."
