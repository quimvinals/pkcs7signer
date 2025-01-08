## Ho to create a signed order

1. [Install Golang](https://go.dev/doc/install) > 1.23
2. At the top level of the repository, add: private key and order type files in .pem extension
    - ordertype_cert.pem
    - private_key.pem
3. At the top level of the repository, add the order you want to sign in a folder called "order"
4. Execute "./generate_signed_order.sh"

## How to verify the signature

Execute ./verify_signature.sh