SERVER_KEY=server.key
SERVER_CRT=server.crt
SERVER_CSR=server.csr
HOST_IP=localhost

setup:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go

# Generate simple .key/.crt using openssl
key:
	openssl genrsa -out ${SERVER_KEY} 2048
	openssl req -new -x509 -sha256 -key ${SERVER_KEY} \
		-out ${SERVER_CRT} -days 3650

# Generate a certificate signing request (.csr) using openssl
csr:
	openssl req -new -sha256 -key ${SERVER_KEY} -out ${SERVER_CSR}
	openssl x509 -req -sha256 -in ${SERVER_CSR} -signkey ${SERVER_KEY} \
		-out ${SERVER_CRT} -days 3650

# Alternative is to use certstrap by Square
# $ brew install certstrap
certstrap:
	# Create a new certificate authority
	certstrap init --common-name alextanhongpin.com
	
	# To request a certificate for a specific host
	certstrap request-cert -ip ${HOST_IP}

	# To generate the certificate for the host:
	certstrap sign ${HOST_IP} --CA alextanhongpin.com

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

serve:
	go run server/main.go

proxy:
	go run gateway/main.go

test: 
	curl -XPOST -d '{"x":1,"y":2}' http://localhost:9090/v1/math/sum