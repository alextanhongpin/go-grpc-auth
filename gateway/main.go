package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"google.golang.org/grpc/credentials"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/alextanhongpin/go-grpc-auth/proto"
)

func run() error {
	var (
		cert = flag.String("cert", "", "path to PEM-formatted CA certificate")
		port = flag.String("port", ":9090", "TCP address to listen on (in host:port) format")
		addr = flag.String("addr", "localhost:8080", "Address of the GRPC service")
	)
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	var opts []grpc.DialOption
	if *cert != "" {
		log.Println("running secure gateway")
		creds, err := credentials.NewClientTLSFromFile(*cert, "")
		if err != nil {
			log.Println(err)
		}
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		log.Println("running insecure gateway")
		opts = []grpc.DialOption{grpc.WithInsecure()}
	}

	err := gw.RegisterMathServiceHandlerFromEndpoint(ctx, mux, *addr, opts)
	if err != nil {
		return err
	}
	log.Printf("listening to port *:%s\n", *port)
	return http.ListenAndServe(*port, mux)
}

func main() {
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
