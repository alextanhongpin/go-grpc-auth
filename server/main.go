package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/alextanhongpin/go-grpc-auth/proto"
)

type mathServer struct{}

var (
	crt = "server.crt"
	key = "server.key"
)

func (s *mathServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	x := req.GetX()
	y := req.GetY()

	z := x + y
	return &pb.SumResponse{
		Z: z,
	}, nil
}

func main() {
	var (
		port = flag.String("port", ":8080", "TCP address to listen on(in host:port form")
		cert = flag.String("cert", "", "Path to PEM-encoded certificate")
		key  = flag.String("key", "", "Path to PEM-encoded secret key")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcServer *grpc.Server

	if *cert == "" && *key == "" {
		log.Println("running insecure server")
		grpcServer = grpc.NewServer()
	} else if *cert == "" {
		log.Fatalf("cert is not available")
	} else if *key == "" {
		log.Fatalf("key is not available")
	} else {
		log.Println("running secure server")
		creds, err := credentials.NewServerTLSFromFile(*cert, *key)
		if err != nil {
			log.Println(err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
	}

	pb.RegisterMathServiceServer(grpcServer, &mathServer{})
	log.Printf("listening to port *:%s\n", *port)
	grpcServer.Serve(lis)

}
