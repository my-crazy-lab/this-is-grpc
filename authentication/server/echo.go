package server

import (
	"context"

	pb "github.com/my-crazy-lab/this-is-grpc/authentication/proto/echo"
)

type ecServer struct {
	pb.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(_ context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: req.Message + " res from auth service"}, nil
}

func NewEchoServer() pb.EchoServer {
	return &ecServer{}
}
