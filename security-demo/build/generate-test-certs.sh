#!/bin/bash 

ret=$(openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) -keyout cert.key -out cert.crt -days 3650 -config ./opensslext.config)

if [ $? -ne 0 ]; then
  echo "Error : dummy certificate generation failed"
  exit 1
else 
  echo "Success : dummy certificate generated"
  exit 0
fi

