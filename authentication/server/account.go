package server

import (
	"context"
	"fmt"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	authPb "github.com/my-crazy-lab/this-is-grpc/authentication/proto/account"
)

type authServer struct {
	authPb.UnimplementedAccountServer
}

func (s *authServer) Login(_ context.Context, req *authPb.AccountRequest) (*authPb.AccountResponse, error) {
	fmt.Printf("innnnn ")
	user, err := pg.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	// Verify password (in production, use hashing)
	if user.Password != req.Password {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT Token
	token, err := pg.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	// Return response with JWT token
	return &authPb.AccountResponse{Token: token}, nil
	// return &authPb.AccountResponse{Token: req.PhoneNumber + " res from auth service"}, nil
}

func NewAuthServer() authPb.AccountServer {
	return &authServer{}
}
