package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	userPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/user"
	"google.golang.org/grpc/metadata"
)

type authServer struct {
	authPb.UnimplementedAuthServer
}

func NewAuthServer() authPb.AuthServer {
	return &authServer{}
}

func verifyToken(ctx context.Context) (*userPb.User, error) {
	// Extract token from gRPC metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authHeader, exists := md["user-authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, errors.New("unauthorized: missing token")
	}

	token := authHeader[0]
	userId, err := pg.VerifyJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return pg.GetUserById(int32(userId))
}
