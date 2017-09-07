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

var cert = "server.crt"

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	// creds := credentials.NewTLS(&tls.Config{
	// 	ServerName: "localhost:8080",
	// 	RootCAs:    "../server.crt",
	// })
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Println(err)
	}
	// opts := []grpc.DialOption{grpc.WithInsecure()}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	err = gw.RegisterMathServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		return err
	}
	log.Printf("listening to port *:%s\n", "8080")
	return http.ListenAndServe(":9090", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
