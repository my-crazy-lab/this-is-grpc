package schema

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/graphql-go/graphql"
	ecpb "github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/proto/echo"
	"github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/rpcServices"
)

type Account struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}
type Auth struct {
	Token string `json:"token"`
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
var TheAuth = Auth{Token: "aaa"}

var accountType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"phoneNumber": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authMutation = graphql.Fields{
	/*
	   curl -g 'http://localhost:9090/graphql?query={signIn{phoneNumber}}'
	*/
	"signIn": &graphql.Field{
		Type:        authType,
		Description: "User sign in",
		Args: graphql.FieldConfigArgument{
			"phoneNumber": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			rgc := rpcServices.NewAuthenticationService()
			callUnaryEcho(rgc, "hello world")

			return TheAuth, nil
		},
	},
}

var authQuery = graphql.Fields{
	/*
	   curl -g 'http://localhost:9090/graphql?query={getUserDetail(token:"token"){userId,username}}'
	*/
	"getUserDetail": &graphql.Field{
		Type:        userType,
		Description: "Get user information",
		Args: graphql.FieldConfigArgument{
			"token": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return User{}, nil
		},
	},
}

func callUnaryEcho(client ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}
