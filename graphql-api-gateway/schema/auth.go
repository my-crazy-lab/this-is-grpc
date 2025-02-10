package schema

import (
	"context"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	authPb "github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/proto/auth"

	"github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/rpcServices"
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
	"signIn": &graphql.Field{
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

			token := login(rpcServices.AuthenticationService, phoneNumber, password)

			return LoginResponse{token}, nil
		},
	},
	"signUp": &graphql.Field{
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
				msg := register(rpcServices.AuthenticationService, phoneNumber, password)
				return msg, nil
			}
			return "", nil
		},
	},
}

var authQuery = graphql.Fields{
	"getUsers": &graphql.Field{
		Type:        graphql.NewList(userType),
		Description: "Get all users",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context.(context.Context)

			users := getUsers(ctx, rpcServices.AuthenticationService)
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

func register(client authPb.AuthClient, phone string, pass string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &authPb.RegisterRequest{PhoneNumber: phone, Password: pass})
	if err != nil {
		log.Fatalf("client.register(_) = _, %v: ", err)
	}

	return resp.Msg
}

func getUsers(ctx context.Context, client authPb.AuthClient) []*authPb.User {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.GetUsers(ctx, &authPb.GetUsersRequest{})
	if err != nil {
		log.Fatalf("client.GetUsers(_) = _, %v: ", err)
	}

	return resp.Users
}
