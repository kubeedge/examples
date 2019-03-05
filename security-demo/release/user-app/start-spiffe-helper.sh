#!/bin/sh

cd certs 
bash -x generate-bundle.sh
cd ..
/opt/spire/spiffe-helper -config /opt/spire/user-app/user-app-helper.conf &> /opt/spire/log/user-app-helper.log &
