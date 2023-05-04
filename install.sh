#!/bin/bash
#
# cx installer
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/liferay/liferay-cli/master/install.sh | bash

# When releasing cx, the releaser should update this version number
# AFTER they upload new binaries.
VERSION="0.1.4"

set -e

function copy_binary() {
  if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
      if [ ! -d "$HOME/.local/bin" ]; then
        mkdir -p "$HOME/.local/bin"
      fi
      mv cx "$HOME/.local/bin/cx"
  elif [[ ":$PATH:" == *":$HOME/bin:"* ]]; then
      if [ ! -d "$HOME/bin" ]; then
        mkdir -p "$HOME/bin"
      fi
      mv cx "$HOME/bin/cx"
  else
      echo "Installing cx to /usr/local/bin which is write protected"
      echo "If you'd prefer to install cx without sudo permissions, add \$HOME/.local/bin OR \$HOME/bin to your \$PATH and rerun the installer"
      sudo mv cx /usr/local/bin/cx
  fi
}

function install_cx() {
  if [[ "$OSTYPE" == "linux"* ]]; then
        # On Linux, "uname -m" reports "aarch64" on ARM 64 bits machines,
        # and armv7l on ARM 32 bits machines like the Raspberry Pi.
        # This is a small workaround so that the install script works on ARM.
        case $(uname -m) in
            aarch64) ARCH=arm64;;
            armv7l)  ARCH=arm;;
            x86_64)  ARCH=amd64;;
            *)       ARCH=$(uname -m);;
        esac
        set -x
        curl -fsSL https://github.com/liferay/liferay-cli/releases/download/v$VERSION/liferay-linux-$ARCH -o cx
        chmod +x cx
        copy_binary
  elif [[ "$OSTYPE" == "darwin"* ]]; then
        # On macOS, "uname -m" reports "arm64" on ARM 64 bits machines
        ARCH=$(uname -m)
        set -x
        curl -fsSL https://github.com/liferay/liferay-cli/releases/download/v$VERSION/liferay-darwin-$ARCH -o cx
        chmod +x cx
        copy_binary
  else
      set +x
      echo "The cx installer does not work for your platform: $OSTYPE"
      echo "For other installation options, check the following page:"
      echo "https://github.com/liferay/liferay-cli#manuall-installation"
      echo "If you think your platform should be supported, please file an issue:"
      echo "https://github.com/liferay/liferay-cli/issues/new"
      echo "Thank you!"
      exit 1
  fi

  set +x
}

function version_check() {
  VERSION_FROM_BIN="$(cx --version 2>&1 || true)"
  CX_DEV_PATTERN='^liferay version [0-9]+\.[0-9]+\.[0-9]+(.*)?$'
  if ! [[ $VERSION_FROM_BIN =~ $CX_DEV_PATTERN ]]; then
    echo "cx installed!"
    echo
    echo "Note: it looks like it is not the first program named 'cx' in your path. \`cx --version\` (running from $(command -v cx)) did not return a cx.dev version string."
    echo "It output this instead:"
    echo
    echo "$VERSION_FROM_BIN"
    echo
    echo "Perhaps you have a different program named cx in your \$PATH?"
    exit 1
  else
    echo "cx installed!"
    echo "Run \`cx ext start\` to start."
  fi
}

# so that we can skip installation in CI and just test the version check
if [[ -z $NO_INSTALL ]]; then
  install_cx
fi

version_check

cx --version