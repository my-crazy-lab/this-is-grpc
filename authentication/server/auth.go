package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	userPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/user"
	"google.golang.org/grpc/metadata"
)

func verifyToken(ctx context.Context) (*pg.User, error) {
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

	return pg.GetUserById(userId)
}

func (s *authServer) Login(_ context.Context, req *authPb.LoginRequest) (*authPb.LoginResponse, error) {
	user, err := pg.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	if !pg.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := pg.GenerateJWT(user.ID)

	if err != nil {
		return nil, err
	}

	return &authPb.LoginResponse{Token: token}, nil
}

func (s *authServer) Register(_ context.Context, req *authPb.RegisterRequest) (*authPb.RegisterResponse, error) {
	// Check if the user already exists
	existingUser, err := pg.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	if existingUser != nil {
		fmt.Println("User already exists:", existingUser)
	}

	err = pg.InsertNewUser(req.PhoneNumber, req.Password)
	if err != nil {
		return nil, err
	}

	return &authPb.RegisterResponse{Msg: "register successful"}, nil
}

func (s *authServer) GetUsers(ctx context.Context, _ *authPb.GetUsersRequest) (*authPb.GetUsersResponse, error) {
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
	_, err := pg.VerifyJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	users, err := pg.GetUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []*userPb.User
	for _, user := range users {
		userResponses = append(userResponses, &userPb.User{
			Id:          int32(user.ID),
			PhoneNumber: user.PhoneNumber,
		})
	}

	return &authPb.GetUsersResponse{Users: userResponses}, nil
}
