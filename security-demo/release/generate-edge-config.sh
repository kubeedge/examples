#!/bin/bash

. ./edge.env

#edge core config
ret=$(sed -i "s#fb4ebb70-2783-42b8-b3ef-63e2fd6d242e#$NODE_ID#" ./app-binaries/edge/conf/edge.yaml)
ret=$(sed -i "s#e632aba927ea4ac2b575ec1603d56f10#$PROJECT_ID#" ./app-binaries/edge/conf/edge.yaml)

ret=$(sed -i "s#1883#$MQTT_EXT_PORT#" ./app-binaries/edge/conf/edge.yaml)
ret=$(sed -i "s#1884#$MQTT_INT_PORT#" ./app-binaries/edge/conf/edge.yaml)

ret=$(sed -i "s#wss://0.0.0.0:10000#wss://$EDGE_HUB_IP:$EDGE_HUB_PORT#" ./app-binaries/edge/conf/edge.yaml)

#spire agent config
ret=$(sed -i "s#bind_address = \"192.168.56.101\"#bind_address = \"$EDGE_VM_IP\"#" ./conf/agent/agent.conf) 
ret=$(sed -i "s#server_address = \"192.168.56.101\"#bind_address = \"$CLOUD_VM_IP\"#" ./conf/agent/agent.conf) 

#app spire agent config
ret=$(sed -i "s#192.168.56.102#$EDGE_VM_IP#" ./app-agent-conf/agent/agent.conf)

#event bus helper
ret=$(sed -i "s#192.168.56.102#$EDGE_VM_IP#" ./event-bus/event-bus-helper.conf)

#user app helper
ret=$(sed -i "s#192.168.56.102#$EDGE_VM_IP#" ./user-app/user-app-helper.conf)

#spire server config
ret=$(sed -i "s#\"disk\"#\"spire\"#" ./conf/server/server.conf)
echo "\
plugins {
    UpstreamCA \"spire\" {
        plugin_data {
            server_address: \"$CLOUD_VM_IP\"
            server_port: \"8081\"
            workload_api_socket: \"/tmp/agent.sock\"
          }
    }
}" >> ./conf/server/server.conf

#cloud hub helper config
ret=$(sed -i "s#127.0.0.1:20000#$CLOUD_VM_IP:40000#" ./helper.conf) 
ret=$(sed -i "s#192.168.56.101:40000#$EDGE_HUB_IP:$EDGE_HUB_PORT#" ./helper.conf) 

