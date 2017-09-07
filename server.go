package main

import (
	"fmt"
)

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

func GrpcServerTLS (certFile, keyFile string) (*grpc.Server, error) {

	var server *grpc.Server
	if keyFile == "" && certFile == "" {
		server = grpc.NewServer()
		return server, nil
	} else if certFile == "" {
		// dieIf(fmt.Errorf("key specified with no cert"))
		return nil, errors.New("key specified with no cert")
	} else if keyFile == "" {
		// dieIf(fmt.Errorf("cert specified with no key"))
		return nil, errors.New("cert specified with no cert")
	} 
	pair, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	creds := grpc.Creds(pair)
	server = grpc.NewServer(creds)

	return server, nil
}

func main () {
	addr := flag.String("srv", ":8181", "TCP address to listen on(in host:port form")
	certFile := flag.String("cert", "", "Path to PEM-encoded certificate")
	keyFile := flag.String("key", "", "Path to PEM-encoded secret key")

	flag.Parse()

	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("Expecting zero arguments, but got %d", flag.NArg()))
	}

	svg := &genSvc{}

	// var server *grpc.Server
	// if *keyFile == "" && *certFile == "" {
	// 	server = grpc.NewServer()
	// } else if *certFile == "" {
	// 	dieIf(fmt.Errorf("key specified with no cert"))
	// } else if *keyFile == "" {
	// 	dieIf(fmt.Errorf("cert specified with no key"))
	// } else {
	// 	pair, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	dieIf(err)
	// 	creds := grpc.Creds(pair)
	// 	server = grpc.NewServer(creds)
	// }

	server := GrpcServerTLS(*certFile, *keyFile)

	list, err := net.Listen("tcp", *addr)
	dieIf(err)
	pb.RegisterGenSvcServer(server, svc)
	server.Serve(lis)
}
