cd certs ; bash -x generate-bundle.sh
cd ..;./spiffe-helper &> ./log/helper.log &
