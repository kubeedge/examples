#!/bin/sh

cd certs 
bash -x generate-bundle.sh
cd ..
/opt/spire/spiffe-helper -config /opt/spire/event-bus/event-bus-helper.conf &> /opt/spire/log/event-bus-helper.log &
