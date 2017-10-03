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

REM Clean errorlevel
verify > nul

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
	REM echo "Getting webrtc ..."
	REM md %WEBRTC_DIR%
	REM pushd %WEBRTC_DIR%
	REM verify > nul
	REM echo "gclient config ..."
	REM call gclient config --name src %WEBRTC_REPO%
    REM IF /I "%ERRORLEVEL%" NEQ "0" exit /b
	REM call gclient sync --with_branch_heads -r %COMMIT%
    REM IF /I "%ERRORLEVEL%" NEQ "0" exit /b
	REM popd
REM )

echo "Checking out latest tested / compatible version of webrtc ..."
pushd %WEBRTC_SRC%
call git checkout %COMMIT%
IF /I "%ERRORLEVEL%" NEQ "0" exit /b
popd

echo "Cleaning webrtc ..."
pushd %WEBRTC_SRC%
IF /I "%ERRORLEVEL%" NEQ "0" exit /b
SET OUTDIR=out\%CONFIG%
echo "Deleting build configs in" %OUTDIR%
rmdir /S /Q %OUTDIR%
popd

echo "Building webrtc ..."
pushd %WEBRTC_SRC%
SET GNARGS=%WEBRTC_SRC%\out\%CONFIG% --args="is_debug=false"
call gn gen %GNARGS%
IF /I "%ERRORLEVEL%" NEQ "0" exit /b

SET NINJAARGS=-C  %WEBRTC_SRC%\out\%CONFIG% webrtc field_trial metrics_default pc_test_utils
call ninja %NINJAARGS%
IF /I "%ERRORLEVEL%" NEQ "0" exit /b
popd

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
)
setlocal

REM pushd $PROJECT_DIR || exit 1
REM git clean -fd "$INCLUDE_DIR"

"C:\Program Files (x86)\Microsoft Visual Studio 14.0\VC\vcvarsall.bat" amd64

echo "Concatenating libraries ..."
pushd %WEBRTC_SRC%\out\%CONFIG%\obj\webrtc
echo %CD%
REM if [ "$OS" = "darwin" ]; then
REM 	find . -name '*.o' > filelist
REM 	libtool -static -o libwebrtc-magic.a -filelist filelist
REM 	strip -S -x -o libwebrtc-magic.a libwebrtc-magic.a
REM elif [ "$ARCH" = "arm" ]; then
REM 	arm-linux-gnueabihf-ar crs libwebrtc-magic.a $(find . -name '*.o' -not -name '*.main.o')
REM else
REM 	ar crs libwebrtc-magic.a $(find . -name '*.o' -not -name '*.main.o')
REM fi
lib /OUT:libwebrtc-magic.lib webrtc.lib webrtc_common.lib
SET OUT_LIBRARY=%LIB_DIR%\libwebrtc-%OS%-%ARCH%-magic.lib
echo %OUT_LIBRARY%
move libwebrtc-magic.lib %OUT_LIBRARY%
echo Built %OUT_LIBRARY%
popd

echo "Build complete."

popd

