package main

import (
	"context"
	"fmt"
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
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", "8080"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		log.Println(err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterMathServiceServer(grpcServer, &mathServer{})
	log.Printf("listening to port *:%s\n", "8080")
	grpcServer.Serve(lis)

}
