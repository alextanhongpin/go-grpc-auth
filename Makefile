SERVER_KEY=server.key
SERVER_CRT=server.crt
SERVER_CSR=server.csr
HOST_IP=172.16.21.215

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
