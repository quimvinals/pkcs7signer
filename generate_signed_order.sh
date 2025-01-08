rm ./order/manifest.json

rm ./order/signature

go mod tidy

go mod vendor

go run manifest/manifest.go ./order ./order/manifest.json

go run signer/signer.go  ./order/manifest.json ordertype_cert.pem private_key.pem ./order/signature