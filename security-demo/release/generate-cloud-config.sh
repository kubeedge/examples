#!/bin/bash

. ./cloud.env

replacestring="kubeconfig: \"\""

#edge controller config
ret=$(sed -i "s#$replacestring#kubeconfig: \"$KUBECONFIG\"#" ./app-binaries/cloud/conf/controller.yaml)
ret=$(sed -i "s#port: 10000#port: $CLOUD_HUB_PORT#" ./app-binaries/cloud/conf/controller.yaml)
ret=$(sed -i "s#address: 0.0.0.0#address: $CLOUD_HUB_IP#" ./app-binaries/cloud/conf/controller.yaml)

#spire agent config
ret=$(sed -i "s#192.168.56.101#$CLOUD_VM_IP#" ./conf/agent/agent.conf) 

#spire server config
if [ ! -f ./conf/server/server.conf.bk ]; then
  cp ./conf/server/server.conf ./conf/server/server.conf.bk
fi

echo "\
plugins {
    UpstreamCA \"disk\" {
        plugin_data {
            ttl = \"1h\"
            key_file_path = \"./conf/server/dummy_upstream_ca.key\"
            cert_file_path = \"./conf/server/dummy_upstream_ca.crt\"
        }
    }
}" >> ./conf/server/server.conf

#cloud hub helper config
ret=$(sed -i "s#127.0.0.1:20000#$CLOUD_HUB_IP:$CLOUD_HUB_PORT#" ./helper.conf) 

