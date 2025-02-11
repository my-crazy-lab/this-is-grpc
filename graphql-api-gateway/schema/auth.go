package schema

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	userPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/user"
	"google.golang.org/grpc/metadata"

	client "github.com/my-crazy-lab/this-is-grpc/proto-module/client"
)

type Account struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type UserType struct {
	Id          string `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

var authType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Auth",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"phoneNumber": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authMutation = graphql.Fields{
	"SignIn": &graphql.Field{
		Type:        authType,
		Description: "User sign in",
		Args: graphql.FieldConfigArgument{
			"phoneNumber": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			phoneNumber, _ := params.Args["phoneNumber"].(string)
			password, _ := params.Args["password"].(string)

			token := login(client.AuthenticationService, phoneNumber, password)

			return LoginResponse{token}, nil
		},
	},
	"SignUp": &graphql.Field{
		Type:        authType,
		Description: "User sign up",
		Args: graphql.FieldConfigArgument{
			"phoneNumber": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			phoneNumber, okPhone := params.Args["phoneNumber"].(string)
			password, okPass := params.Args["password"].(string)

			if okPhone && okPass {
				msg := register(client.AuthenticationService, &authPb.RegisterRequest{PhoneNumber: phoneNumber, Password: password})
				return msg, nil
			}
			return "", nil
		},
	},
}

var authQuery = graphql.Fields{
	"GetUsers": &graphql.Field{
		Type:        graphql.NewList(userType),
		Description: "Get all users",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			users, err := getUsers(ctx, client.AuthenticationService)
			if err != nil {
				log.Fatalf("get users error %v: ", err)
			}
			return users, nil
		},
	},
	"GetUser": &graphql.Field{
		Type:        userType,
		Description: "Get user by id",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			users, err := getUser(ctx, client.AuthenticationService)
			if err != nil {
				log.Fatalf("get users error %v: ", err)
			}
			return users, nil
		},
	},
}

func login(client authPb.AuthClient, phone string, pass string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &authPb.LoginRequest{PhoneNumber: phone, Password: pass})
	if err != nil {
		log.Fatalf("client.login(_) = _, %v: ", err)
	}

	return resp.Token
}

func register(client authPb.AuthClient, params *authPb.RegisterRequest) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, params)
	if err != nil {
		log.Fatalf("client.register(_) = _, %v: ", err)
	}

	return resp.Msg
}

func getUsers(ctx context.Context, client authPb.AuthClient) ([]*userPb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetUsers(ctx, &authPb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("client.GetUsers(_) = _, %v: ", err)
	}

	return resp.Users, nil
}

func getUser(ctx context.Context, client authPb.AuthClient) (*userPb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetUser(ctx, &authPb.GetUserRequest{})
	if err != nil {
		log.Fatalf("client.GetUser(_) = _, %v: ", err)
	}

	return resp, nil
}
