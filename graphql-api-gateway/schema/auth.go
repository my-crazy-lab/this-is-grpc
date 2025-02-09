package schema

import (
	"context"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	authPb "github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/proto/account"

	"github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/rpcServices"
)

type Account struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

var authType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Auth",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var accountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"phoneNumber": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
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
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			phoneNumber, _ := params.Args["phoneNumber"].(string)

			token := login(rpcServices.AuthenticationService, phoneNumber, "")

			return token, nil
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
			phoneNumber, _ := params.Args["phoneNumber"].(string)
			password, _ := params.Args["password"].(string)

			msg := register(rpcServices.AuthenticationService, phoneNumber, password)
			return msg, nil
		},
	},
}

var authQuery = graphql.Fields{
	"getUsers": &graphql.Field{
		Type:        graphql.NewList(accountType),
		Description: "Get all users",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			users := getUsers(rpcServices.AuthenticationService)
			return users, nil
		},
	},
}

func login(client authPb.AccountClient, phone string, pass string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Login(ctx, &authPb.AccountRequest{PhoneNumber: phone, Password: pass})
	if err != nil {
		log.Fatalf("client.login(_) = _, %v: ", err)
	}

	return resp.Token
}

func register(client authPb.AccountClient, phone string, pass string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &authPb.AccountRequest{PhoneNumber: phone, Password: pass})
	if err != nil {
		log.Fatalf("client.register(_) = _, %v: ", err)
	}

	return resp.Msg
}

func getUsers(client authPb.AccountClient) []*authPb.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetUsers(ctx, &authPb.Empty{})
	if err != nil {
		log.Fatalf("client.GetUsers(_) = _, %v: ", err)
	}

	return resp.Users
}
