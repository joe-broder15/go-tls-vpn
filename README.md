# go-tls-vpn
TLS vpn written in golang

## Certificates
The server will require an x509 key pair in order to work properly. You can generate a key pair via the following commands
```
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
