package server

import authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"

type authServer struct {
	authPb.UnimplementedAuthServer
}

func NewAuthServer() authPb.AuthServer {
	return &authServer{}
}
