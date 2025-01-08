openssl pkcs7 -inform der -in ./order/signature -out signature

openssl pkcs7 -print_certs -in signature -out signature.cert 

openssl smime -verify -binary -inform PEM -in signature -content ./order/manifest.json -certfile signature.cert -nointern -noverify > /dev/null

rm signature signature.cert