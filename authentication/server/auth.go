package server

import (
	"context"
	"fmt"
	"log"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	userPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/user"
)

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
	verifyToken(ctx)

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

func (s *authServer) GetUser(ctx context.Context, req *authPb.GetUserRequest) (*userPb.User, error) {
	verifyToken(ctx)

	user, err := pg.GetUserById(req.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
