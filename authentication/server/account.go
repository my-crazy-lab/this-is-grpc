package server

import (
	"context"
	"fmt"
	"log"

	"github.com/my-crazy-lab/this-is-grpc/authentication/pg"
	authPb "github.com/my-crazy-lab/this-is-grpc/authentication/proto/account"
)

type authServer struct {
	authPb.UnimplementedAccountServer
}

func (s *authServer) Login(_ context.Context, req *authPb.AccountRequest) (*authPb.AccountResponse, error) {
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

	return &authPb.AccountResponse{Token: token}, nil
}

func (s *authServer) Register(_ context.Context, req *authPb.AccountRequest) (*authPb.Msg, error) {
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

	return &authPb.Msg{Msg: "register successful"}, nil
}

func (s *authServer) GetUsers(ctx context.Context, _ *authPb.Empty) (*authPb.UsersResponse, error) {
	token, ok := ctx.Value("token").(string)
	if !ok {
		return nil, fmt.Errorf("unauthorized access: token not found in context")
	}

	_, err := pg.VerifyJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	users, err := pg.GetUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []*authPb.User
	for _, user := range users {
		userResponses = append(userResponses, &authPb.User{
			Id:          int32(user.ID),
			PhoneNumber: user.PhoneNumber,
			Password:    user.Password,
		})
	}

	return &authPb.UsersResponse{Users: userResponses}, nil
}

func NewAuthServer() authPb.AccountServer {
	return &authServer{}
}
