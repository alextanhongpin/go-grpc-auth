# go-grpc-auth

This is a simple grpc demo with auth.

## Setup

You need to have `golang` installed. Run `make setup` to install other dependencies:

```Makefile
setup:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go
```

## GRPC Stub and Reverse-Proxy

Run `make stub` and `make gw`:

```Makefile
stub:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. \
	proto/*.proto

gw:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	proto/*.proto
```

This will generate `hello.pb.go` && `hello.pb.gw.go` respectively from your `hello.proto` file.

## Cert and Key

To generate the cert and key, run `make key` and `make csr`:

```Makefile
key:
	openssl genrsa -out ${SERVER_KEY} 2048
	openssl req -new -x509 -sha256 -key ${SERVER_KEY} \
		-out ${SERVER_CRT} -days 3650

csr:
	openssl req -new -sha256 -key ${SERVER_KEY} -out ${SERVER_CSR}
	openssl x509 -req -sha256 -in ${SERVER_CSR} -signkey ${SERVER_KEY} \
		-out ${SERVER_CRT} -days 3650
```

Make sure the field `Common Name (e.g. server FQDN or YOUR name) []:` is filled with your server hostname. It's `localhost` if you run it locally.

```bash
$ make key

openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus
.................+++
..................................................................................................................................................................+++
e is 65537 (0x10001)
openssl req -new -x509 -sha256 -key server.key \
		-out server.crt -days 3650
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:.
State or Province Name (full name) [Some-State]:.
Locality Name (eg, city) []:.
Organization Name (eg, company) [Internet Widgits Pty Ltd]:.
Organizational Unit Name (eg, section) []:.
Common Name (e.g. server FQDN or YOUR name) []:localhost
Email Address []:.
```

```bash
$ make csr

openssl req -new -sha256 -key server.key -out server.csr
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:.
State or Province Name (full name) [Some-State]:.
Locality Name (eg, city) []:.
Organization Name (eg, company) [Internet Widgits Pty Ltd]:.
Organizational Unit Name (eg, section) []:.
Common Name (e.g. server FQDN or YOUR name) []:localhost
Email Address []:.

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:.
An optional company name []:.
openssl x509 -req -sha256 -in server.csr -signkey server.key \
		-out server.crt -days 3650
Signature ok
subject=/CN=localhost
Getting Private key
```

## Run Server

Now that you have both the cert and key, run `make sserve` and `make sproxy`. `sserve` stands for *secure-serve*,
`ssproxy` for *secure-proxy*. 

To test the output, run `make test`. You should get the following output:

```bash
curl -XPOST -d '{"x":1,"y":2}' http://localhost:9090/v1/math/sum
{"z":3}%
```

## Insecure server/gateway

If the gateway is insecure and you try to make a call through `make test`, you should get the following error:

```bash
$ curl -XPOST -d '{"x":1,"y":2}' http://localhost:9090/v1/math/sum
{"error":"transport is closing","code":14}%
```