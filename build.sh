#!/bin/bash
# Initial dependency script ensures libwebrtc is available
# locally.
PROJECT_DIR=$(pwd)
LIBWEBRTC_REPO="https://github.com/js-platform/libwebrtc.git"
LIB_DIR="$PROJECT_DIR/third_party/"
LIBWEBRTC_DIR=$LIB_DIR"libwebrtc"

echo "libwebrtc lives at: $LIBWEBRTC_DIR"

if [[ -d $LIBWEBRTC_DIR ]]; then
  echo "Updating to latest"
  cd $LIBWEBRTC_DIR
  git pull origin master
  cd ..
else
  echo "Cloning new libwebrtc"
  git clone $LIBWEBRTC_REPO $LIBWEBRTC_DIR
fi
