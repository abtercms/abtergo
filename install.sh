#!/usr/bin/env bash

# install task
mkdir -p build
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.bin

echo "Run task install to install additional tools"
