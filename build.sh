#!/usr/bin/env bash

PROJECT_DIR=$(pwd)
THIRD_PARTY_DIR="$PROJECT_DIR/third_party"
WEBRTC_REPO="https://chromium.googlesource.com/external/webrtc"
WEBRTC_DIR="$THIRD_PARTY_DIR/webrtc"
WEBRTC_SRC="$WEBRTC_DIR/src"
DEPOT_TOOLS_DIR="$THIRD_PARTY_DIR/depot_tools"
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
CONFIG="Release"
COMMIT="c279861207c5b15fc51069e96595782350e0ac12"  # branch-heads/58

# Values are from,
#   https://github.com/golang/go/blob/master/src/go/build/syslist.go
#   https://chromium.googlesource.com/chromium/src/+/master/tools/gn/docs/reference.md

oses=",linux:linux,darwin:mac,windows:win,android:android,"
cpus=",386:x86,amd64:x64,arm:arm,"

get() {
	echo "$(expr "$1" : ".*,$2:\([^,]*\),.*")"
}

TARGET_OS=$(get $oses $OS)
TARGET_CPU=$(get $cpus $ARCH)

INCLUDE_DIR="$PROJECT_DIR/include"
LIB_DIR="$PROJECT_DIR/lib"

PATH="$PATH:$DEPOT_TOOLS_DIR"

mkdir -p $THIRD_PARTY_DIR
mkdir -p $INCLUDE_DIR
mkdir -p $LIB_DIR

if [[ -d $DEPOT_TOOLS_DIR ]]; then
	echo "Syncing depot_tools ..."
	pushd $DEPOT_TOOLS_DIR
	git pull --rebase || exit 1
	popd
else
	echo "Getting depot_tools ..."
	mkdir -p $DEPOT_TOOLS_DIR
	git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git $DEPOT_TOOLS_DIR || exit 1
fi

if [[ -d $WEBRTC_DIR ]]; then
	echo "Syncing webrtc ..."
	pushd $WEBRTC_SRC || exit 1
	if ! git diff-index --quiet HEAD --; then
		echo -en "\nOpen files present in $WEBRTC_SRC\nReset them? (y/N): "
		read ANSWER
		if [ "$ANSWER" != "y" ]; then
			echo "*** Cancelled ***"
			exit 1
		fi
		git reset --hard HEAD || exit 1
	fi
	popd

	pushd $WEBRTC_DIR
	gclient sync --with_branch_heads -r $COMMIT || exit 1
	popd
else
	echo "Getting webrtc ..."
	mkdir -p $WEBRTC_DIR
	pushd $WEBRTC_DIR
	gclient config --name src $WEBRTC_REPO || exit 1
	gclient sync --with_branch_heads -r $COMMIT || exit 1
	popd
fi

if [ "$ARCH" = "arm" ]; then
	echo "Manually fetching arm sysroot"
	pushd $WEBRTC_SRC || exit 1
	./build/linux/sysroot_scripts/install-sysroot.py --arch=arm || exit 1
	popd
fi

echo "Checking out latest tested / compatible version of webrtc ..."
pushd $WEBRTC_SRC
git checkout $COMMIT
popd

echo "Cleaning webrtc ..."
pushd $WEBRTC_SRC || exit 1
rm -rf out/$CONFIG
popd

echo "Building webrtc ..."
pushd $WEBRTC_SRC
gn gen out/$CONFIG --args="target_os=\"$TARGET_OS\" target_cpu=\"$TARGET_CPU\" is_debug=false" || exit 1
ninja -C out/$CONFIG webrtc field_trial metrics_default pc_test_utils || exit 1
popd

echo "Copying headers ..."
pushd $WEBRTC_SRC || exit 1
rm -rf "$INCLUDE_DIR"
for h in $(find webrtc/ -type f -name '*.h')
do
	mkdir -p "$INCLUDE_DIR/$(dirname $h)"
	cp $h "$INCLUDE_DIR/$h"
done
popd
pushd $PROJECT_DIR || exit 1
git clean -fd "$INCLUDE_DIR"
popd

echo "Concatenating libraries ..."
pushd $WEBRTC_SRC/out/$CONFIG
if [ "$OS" = "darwin" ]; then
	find obj -name '*.o' > filelist
	libtool -static -o libwebrtc-magic.a -filelist filelist
	strip -S -x -o libwebrtc-magic.a libwebrtc-magic.a
elif [ "$ARCH" = "arm" ]; then
	arm-linux-gnueabihf-ar crs libwebrtc-magic.a $(find obj -name '*.o')
else
	ar crs libwebrtc-magic.a $(find obj -name '*.o')
fi
OUT_LIBRARY=$LIB_DIR/libwebrtc-$OS-$ARCH-magic.a
mv libwebrtc-magic.a ${OUT_LIBRARY}
echo "Built ${OUT_LIBRARY}"
popd

echo "Build complete."
