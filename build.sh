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
COMMIT="f33698296719f956497d2dbff81b5080864a8804"  # branch-heads/52

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

echo "Checking out latest tested / compatible version of webrtc ..."
pushd $WEBRTC_SRC
git checkout $COMMIT
popd

echo "Cleaning webrtc ..."
pushd $WEBRTC_SRC || exit 1
rm -rf out/$CONFIG
popd

echo "Applying webrtc patches ..."
pushd $WEBRTC_SRC || exit 1
for PATCH in build_at_webrtc_branch_heads_52.patch; do
	git apply --check ${PROJECT_DIR}/webrtc_patches/${PATCH} || exit 1
	git am < ${PROJECT_DIR}/webrtc_patches/${PATCH} || exit 1
done
popd

echo "Building webrtc ..."
pushd $WEBRTC_SRC
export GYP_DEFINES="include_tests=0 include_examples=0"
python webrtc/build/gyp_webrtc webrtc/api/api.gyp || exit 1
ninja -C out/$CONFIG || exit 1
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
git clean -fdx "$INCLUDE_DIR"
popd

echo "Concatenating libraries ..."
pushd $WEBRTC_SRC/out/$CONFIG
if [ "$OS" = "darwin" ]; then
	ls *.a > filelist
	libtool -static -o libwebrtc-magic.a -filelist filelist
	strip -S -x -o libwebrtc-magic.a libwebrtc-magic.a
else
	ar crs libwebrtc-magic.a $(find . -name '*.o' -not -name '*.main.o')
fi
OUT_LIBRARY=$LIB_DIR/libwebrtc-$OS-$ARCH-magic.a
mv libwebrtc-magic.a ${OUT_LIBRARY}
echo "Built ${OUT_LIBRARY}"
popd

echo "Build complete."
