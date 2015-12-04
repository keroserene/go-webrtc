#!/bin/bash

PROJECT_DIR=$(pwd)
THIRD_PARTY_DIR="$PROJECT_DIR/third_party"
WEBRTC_REPO="https://chromium.googlesource.com/external/webrtc"
WEBRTC_DIR="$THIRD_PARTY_DIR/webrtc"
WEBRTC_SRC="$WEBRTC_DIR/src"
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
CONFIG="Debug"
COMMIT="cb3f9bd"


# TODO(arlolra): depot_tools

mkdir -p $THIRD_PARTY_DIR

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

echo "Build complete."
