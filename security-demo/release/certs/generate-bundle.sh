rm new_bundle.p12
openssl pkcs12 -export -in svid.pem -inkey svid_key.pem -certfile svid_bundle.pem -out new_bundle.p12
