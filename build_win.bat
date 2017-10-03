echo off
SET PROJECT_DIR=%CD%
SET THIRD_PARTY_DIR=%PROJECT_DIR%\third_party
SET WEBRTC_REPO="https://chromium.googlesource.com/external/webrtc"
SET WEBRTC_DIR=%THIRD_PARTY_DIR%\webrtc
SET WEBRTC_SRC=%WEBRTC_DIR%\src
SET DEPOT_TOOLS_DIR=%THIRD_PARTY_DIR%\depot_tools
SET OS=windows
SET ARCH=amd64
SET CONFIG=Release

REM branch-heads/62
SET COMMIT="6f21dc245689c29730002da09534a8d275e6aa92"

SET TARGET_OS=windows:win
SET TARGET_CPU=amd64:x64

SET INCLUDE_DIR=%PROJECT_DIR%\include\webrtc
SET LIB_DIR=%PROJECT_DIR%\lib

md %THIRD_PARTY_DIR%
md %INCLUDE_DIR%
md %LIB_DIR%

REM IF EXIST %DEPOT_TOOLS_DIR% (
REM 	echo "Syncing depot_tools ..."
REM 	pushd %DEPOT_TOOLS_DIR%
REM 	git pull --rebase
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM 	popd
REM ) ELSE (
REM 	echo "Getting depot_tools ..."
REM 	mkdir -p %DEPOT_TOOLS_DIR%
REM 	git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git %DEPOT_TOOLS_DIR%
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM )

REM IF EXIST %WEBRTC_DIR% (
REM 	echo "Syncing webrtc ..."
REM 	pushd %WEBRTC_SRC%
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM     REM ERS - not sure what this code does, and how to test it
REM 	REM if ! git diff-index --quiet HEAD --; then
REM 	REM 	echo -en "\nOpen files present in $WEBRTC_SRC\nReset them? (y/N): "
REM 	REM 	read ANSWER
REM 	REM 	if [ "$ANSWER" != "y" ]; then
REM 	REM 		echo "*** Cancelled ***"
REM 	REM 		exit 1
REM 	REM 	fi
REM     echo "Issuing a git reset --hard"
REM 	git reset --hard HEAD
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM 	REM fi
REM 	popd

REM 	pushd %WEBRTC_DIR%
REM     echo "Syncing code base to" %COMMIT% "in dir" %WEBRTC_DIR%
REM 	gclient sync --with_branch_heads -r %COMMIT%
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM 	popd
REM ) ELSE (
REM 	echo "Getting webrtc ..."
REM 	md %WEBRTC_DIR%
REM 	pushd $WEBRTC_DIR
REM 	gclient config --name src %WEBRTC_REPO%
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM 	gclient sync --with_branch_heads -r %COMMIT%
REM     IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM 	popd
REM )

REM gclient config --name src https://chromium.googlesource.com/external/webrtc
REM gclient sync --with_branch_heads -r 6f21dc245689c29730002da09534a8d275e6aa92

REM echo "Checking out latest tested / compatible version of webrtc ..."
REM pushd %WEBRTC_SRC%
REM git checkout %COMMIT%
REM IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM popd

REM echo "Cleaning webrtc ..."
REM pushd %WEBRTC_SRC%
REM IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM SET OUTDIR=out\%CONFIG%
REM echo "Deleting build configs in" %OUTDIR%
REM rmdir /S /Q %OUTDIR%
REM popd

REM echo "Building webrtc ..."
REM pushd %WEBRTC_SRC%
REM SET GNARGS=%WEBRTC_SRC%\out\%CONFIG% --args="is_debug=false"
REM call gn gen %GNARGS%
REM IF /I "%ERRORLEVEL%" NEQ "0" exit /b

REM SET NINJAARGS=-C  %WEBRTC_SRC%\out\%CONFIG% webrtc field_trial metrics_default pc_test_utils
REM ninja %NINJAARGS%
REM IF /I "%ERRORLEVEL%" NEQ "0" exit /b
REM popd

echo "Copying headers ..."
pushd %WEBRTC_SRC%\webrtc
IF /I "%ERRORLEVEL%" NEQ "0" exit /b
rmdir /S /Q %INCLUDE_DIR%
setlocal enabledelayedexpansion
for /R . %%f in (*.h) do (
    SET B=%%~dpf
    SET REL=!B:%CD%\=!
    SET DST=%INCLUDE_DIR%\!REL!
    REM echo !DST!
    IF Not EXIST !DST! (
       echo !DST! does not exist
       md !DST!
    )
    SET B=%%f
    SET REL=!B:%CD%\=!
    echo !REL!
    xcopy /S /Y /I !REL! !DST!
REM   SET B=%%f
REM   SET RELATIVE=!B:%CD%\=!
)
REM for h in $(find webrtc/ -type f -name '*.h')
REM do
REM 	mkdir -p "$INCLUDE_DIR/$(dirname $h)"
REM 	cp $h "$INCLUDE_DIR/$h"
REM done
REM popd
REM pushd $PROJECT_DIR || exit 1
REM git clean -fd "$INCLUDE_DIR"
popd

