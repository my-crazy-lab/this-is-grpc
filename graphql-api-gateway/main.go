package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	client "github.com/my-crazy-lab/this-is-grpc/proto-module/client"

	"github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/schema"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operationName"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	client.NewAuthenticationClient()
	defer client.AuthClientConnection.Close()

	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p postData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			w.WriteHeader(400)
			return
		}

		graphqlSchema, err := graphql.NewSchema(graphql.SchemaConfig{
			Query:    schema.RootQuery,
			Mutation: schema.RootMutation,
		})
		if err != nil {
			fmt.Printf("could not declare graphql schema: %s", err)
		}

		ctx := req.Context()
		// Extract the Authorization header
		authHeader := req.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				ctx = context.WithValue(ctx, "token", token)
			}
		}

		result := graphql.Do(graphql.Params{
			Context:        ctx,
			Schema:         graphqlSchema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.Operation,
		})
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("could not write result to response: %s", err)
		}
	})

	fmt.Println("Now server is running on port 9090")

	http.ListenAndServe(":9090", nil)
}
