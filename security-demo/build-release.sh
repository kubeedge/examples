#!/bin/bash

build_dependencies()
{
  echo "Building all dependent libraries.."
  cd build/
  ./generate-test-certs.sh
  if [ $? -ne 0 ]; then
    echo "Release build failed."
    exit 1
  fi
  ./build.sh
  if [ $? -ne 0 ]; then
    echo "Release build failed."
    exit 1
  fi
  cd ..
}

copy_binaries()
{
  echo "Create log directory..."
  mkdir -p release/log

  echo "Copying test certificates..."
  cp build/cert.crt release/app-agent-conf/agent/dummy_root_ca.crt
  cp build/cert.crt release/conf/agent/dummy_root_ca.crt
  cp build/cert.crt release/conf/server/dummy_upstream_ca.crt
  cp build/cert.key release/conf/server/dummy_upstream_ca.key

  echo"Copying configuration files..."
  cp $GOPATH/src/github.com/kubeedge/kubeedge/cloud/edgecontroller/conf/controller.yaml app-binaries/cloud/conf/controller.yaml
  cp $GOPATH/src/github.com/kubeedge/kubeedge/edge/conf/edge.yaml app-binaries/edge/conf/edge.yaml

  echo "Copying libraries..."
  cp $GOPATH/src/github.com/kubeedge/examples/security-demo/cloud-stub/cmd/cloud-app release/app-binaries/cloud/cloud-app && \
  cp $GOPATH/src/github.com/kubeedge/examples/led-raspberrypi/light_mapper/light_mapper release/app-binaries/edge/light_mapper && \
  cp $GOPATH/src/github.com/kubeedge/kubeedge/cloud/edgecontroller/edgecontroller release/app-binaries/cloud/edgecontroller && \
  cp $GOPATH/src/github.com/kubeedge/kubeedge/edge/edge_core release/app-binaries/edge/edge_core && \
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
