#!/bin/bash

build_dependencies()
{
  echo "Building all dependent libraries.."
  cd build/
  ./build.sh
  if [ $? -ne 0 ]; then
    echo "Release build failed."
    exit 1
  fi
  cd ..
}

copy_binaries()
{
  echo "Copying libraries..."
  cp $GOPATH/src/github.com/kubeedge/examples/security-demo/cloud-stub/cmd/cloud-app release/app-binaries/cloud-app && \
  cp $GOPATH/src/github.com/kubeedge/examples/led-raspberrypi/light_mapper/light_mapper release/app-binaries/light_mapper && \
  cp $GOPATH/src/github.com/kubeedge/kubeedge/edge/edge_core release/app-binaries/edge_core && \
  cp $GOPATH/src/github.com/spiffe/spire/cmd/spire-agent/spire-agent release/spire-agent && \
  cp $GOPATH/src/github.com/spiffe/spire/cmd/spire-server/spire-server release/spire-server && \
  cp $GOPATH/src/github.com/spiffe/spiffe-helper/spiffe-helper release/spiffe-helper && \
  cp $GOPATH/src/github.com/spiffe/ghostunnel/ghostunnel release/ghostunnel
  if [ $? -ne 0 ]; then
    echo "Release build failed."
    exit 1
  fi
}

build_dependencies
copy_binaries

echo "Please follow README and use the release folder."
