#!/usr/bin/env bash

# For a native compile (using current values of GOOS and GOARCH):
#   ./build.sh
# For a cross compile:
#   GOOS=linux GOARCH=amd64 ./build.sh
#   GOOS=linux GOARCH=arm ./build.sh
#   GOOS=android GOARCH=arm ./build.sh
# For macOS (GOOS=darwin GOARCH=amd64), you can currently only do a native compile.
# For a cross-compile to linux-arm, you need to install the binutils-arm-linux-gnueabihf package.
# For a cross-compile to android-arm, first run third_party/webrtc/src/build/install-build-deps-android.sh to install needed dependencies.

PROJECT_DIR=$(pwd)
THIRD_PARTY_DIR="$PROJECT_DIR/third_party"
WEBRTC_REPO="https://chromium.googlesource.com/external/webrtc"
WEBRTC_DIR="$THIRD_PARTY_DIR/webrtc"
WEBRTC_SRC="$WEBRTC_DIR/src"
DEPOT_TOOLS_DIR="$THIRD_PARTY_DIR/depot_tools"
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
CONFIG="Release"
COMMIT="88f5d9180eae78a6162cccd78850ff416eb82483"  # branch-heads/64

# Values are from,
#   https://github.com/golang/go/blob/master/src/go/build/syslist.go
#   https://gn.googlesource.com/gn/+/master/docs/reference.md
#
# Android steps from:
#   https://www.chromium.org/developers/gn-build-configuration
#   https://chromium.googlesource.com/chromium/src/+/master/docs/android_build_instructions.md

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

if [[ ! -d $WEBRTC_DIR ]]; then
	echo "Getting webrtc ..."
	mkdir -p $WEBRTC_DIR
	pushd $WEBRTC_DIR
	gclient config --name src $WEBRTC_REPO || exit 1
	popd
fi

if [[ $TARGET_OS == 'android' ]]; then
	echo "Setting gclient target_os to android"
	# Whacky sed to append 'android' to the target_os list, without clobbering what may be there already.
	sed -i "/^target_os *= *.*\\<android\\>/{p;h;d}; /^target_os *= */{s/ *]/, 'android'&/;p;h;d}; \${x;s/^target_os/&/;tx;atarget_os = [ 'android' ]"$'\n'";:x x}" $WEBRTC_DIR/.gclient
fi

echo "Syncing webrtc ..."
if [[ -d $WEBRTC_SRC ]]; then
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
fi
pushd $WEBRTC_DIR
# "echo n" is to say "no" to the Google Play services license agreement and download.
echo n | gclient sync --with_branch_heads -r $COMMIT || exit 1
# Delete where the Google Play services downloads to, just to be sure.
# First check that an ancestor directory of what we're deleting exists, so we're more likely to notice a source reorganization.
if [[ $TARGET_OS == 'android' && ! -d "$WEBRTC_SRC/third_party/android_tools/sdk/extras/google/m2repository" ]]; then
	echo "Didn't find Google Play services directory for removal, please check" 1>&2
	exit 1
fi
rm -rf "$WEBRTC_SRC/third_party/android_tools/sdk/extras/google/m2repository/com/google/android/gms"
popd

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
gn gen out/$CONFIG --args="target_os=\"$TARGET_OS\" target_cpu=\"$TARGET_CPU\" is_debug=false symbol_level=0 use_custom_libcxx=false" || exit 1
ninja -C out/$CONFIG webrtc field_trial metrics_default pc_test_utils || exit 1
popd

echo "Copying headers ..."
pushd $WEBRTC_SRC || exit 1
rm -rf "$INCLUDE_DIR"
find . -type f -name '*.h' -print0 | while IFS= read -r -d '' h;
do
	mkdir -p "$INCLUDE_DIR/$(dirname "$h")"
	cp $h "$INCLUDE_DIR/$h"
done
popd
pushd $PROJECT_DIR || exit 1
git clean -fd "$INCLUDE_DIR"
popd

echo "Concatenating libraries ..."
pushd $WEBRTC_SRC/out/$CONFIG
if [ "$OS" = "darwin" ]; then
	find obj -name '*.o' -print0 \
		| xargs -0 -- libtool -static -o libwebrtc-magic.a
	strip -S -x -o libwebrtc-magic.a libwebrtc-magic.a
elif [ "$OS" = "android" ]; then
	find obj -name '*.o' -print0 | sort -z \
		| xargs -0 -- $WEBRTC_SRC/third_party/android_tools/ndk/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-ar crsD libwebrtc-magic.a
elif [ "$ARCH" = "arm" ]; then
	find obj -name '*.o' -print0 | sort -z \
		| xargs -0 -- arm-linux-gnueabihf-ar crsD libwebrtc-magic.a
else
	find obj -name '*.o' -print0 | sort -z \
		| xargs -0 -- ar crsD libwebrtc-magic.a
fi
OUT_LIBRARY=$LIB_DIR/libwebrtc-$OS-$ARCH-magic.a
mv libwebrtc-magic.a ${OUT_LIBRARY}
echo "Built ${OUT_LIBRARY}"
popd

echo "Build complete."
