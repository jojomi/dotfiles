#!/usr/bin/env sh
#
# For using UPX on linux, goupx (https://github.com/pwaller/goupx) needs to be
# available in $PATH.

UPX_ENABLE=yes
UPX_PARAMS=-9
UPX_PARAMS=--ultra-brute

BIN_DIR="bin"
BIN_LINUX=$BIN_DIR/dotfiles.linux64
BIN_MAC=$BIN_DIR/dotfiles.macosx
BIN_WIN=$BIN_DIR/dotfiles.exe

mkdir --parents "$BIN_DIR"
GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o "$BIN_LINUX" *.go
GOOS=darwin GOARCH=386 go build -ldflags "-w -s" -o "$BIN_MAC" *.go

# minimize binaries using upx (not necessary, might increase number of "security"
# software false positives)
if [ $UPX_ENABLE = "yes" ]; then
  if hash upx 2>/dev/null; then
    upx $UPX_PARAMS "$BIN_MAC"

    GOUPX_BIN=$(which goupx)
    if [ -f $GOUPX_BIN ]; then
      $GOUPX_BIN "$BIN_LINUX" --no-upx
      upx $UPX_PARAMS "$BIN_LINUX"
    fi
  fi
fi
