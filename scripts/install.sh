#!/bin/bash

################################################################################
#                                                                              #
#                                pushx Installer                                #
#                                                                              #
################################################################################

# This script will install the latest version of pushx on your machine from
# the precompiled binary releases in the official repository.

# check platform
if [[ "$(uname)" == "Darwin" ]]; then
  PLATFORM="darwin"
elif [[ "$(uname)" == "Linux" ]]; then
  PLATFORM="linux"
else
  PLATFORM="windows"
fi

BINARIES=${BINARIES:-"pushx"}
INSTALL_DIR=${INSTALL_DIR:-"$HOME/bin"}

install_dir() {
  if [[ -d "$INSTALL_DIR" ]]; then
    echo "$INSTALL_DIR"
  else
    echo /usr/local/bin
  fi
}

install_bin() {
  local name=$1
  if [[ -z "$name" ]]; then
    echo "install_bin: name is empty"
    return 1
  fi
  echo "install_bin: installing $name"
  LATEST_DOWNLOAD_PREFIX="https://github.com/robertlestak/pushx/releases/latest/download/"
  FILE_NAME="${name}_${PLATFORM}"
  DL="${LATEST_DOWNLOAD_PREFIX}${FILE_NAME}"
  echo "install_bin: downloading $DL"
  curl -s -L $DL > $FILE_NAME
  chmod +x $FILE_NAME
  mv $FILE_NAME $(install_dir)/$name
}

install_binaries() {
  # split on comma
  for binary in $(echo $BINARIES | tr "," "\n"); do
    install_bin $binary
  done
}

main() {
  install_binaries
}
main "$@"
