package main

import (
	"flag"
)

func main () {
	caCertFile := flag.String("cacert", "", "path to PEM-formatted CA certificate")
	addr := flag.String("srv", ":8080", "TCP address to listen on (in host:port) format")
	svc := flag.String("grpc-svc", "localhost:8181", "Address of the GRPC service")

	flag.Parse()

	if flag.NArg() != 0 {
		log.Println(flag.NArg())
	}

	var creds grpc.DialOption

	if *caCertFile == "" {
		creds = grpc.WithInsecure()
	} else {
		c, err := crendentials.NewClientTLSFromFile(*caCertFile, *svc)
		if err != nil {
			log.Fatal(err)
		}
		creds = grpc.WithTransportCredentials(c)
	}

	conn, err := grpc.Dial(*svc, creds)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	

}