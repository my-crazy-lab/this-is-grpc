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

	fmt.Println("")

	fmt.Println(`GET USERS:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetUsers { id phoneNumber } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`GET USER BY ID:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetUser { id phoneNumber } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`SIGN IN:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "mutation { SignIn(phoneNumber:\"123456\",password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`REGISTER:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "mutation { SignUp(phoneNumber:\"123456\", password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`CREATE CATEGORIES:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "mutation { CreateCategories(name:\"must unique 2\", description:\"description hihihi\") { id } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`GET CATEGORIES:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetCategories { id  name description created_at updated_at } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`CREATE PRODUCT:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "mutation { CreateProduct(name:\"product 1\", description:\"description hihihi\", price: 200.5, quantity: 2, categories: [1,2]) { id } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`GET PRODUCTS:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetProducts(pageSize:10, pageIndex: 1) { total products {id  name description price quantity categories {id name description created_at updated_at} created_at updated_at} } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`GET PRODUCT BY ID:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJleHAiOjE3MzkyNjY3MTV9.9bIhPrDnO8k7h0gnnZHF2afh7fRAyrwuOz14gdWf8PA" \
--data '{ "query": "{ GetProduct(id: 1) {id  name description price quantity categories {id name description created_at updated_at} created_at updated_at} }" }' \
http://localhost:9090/graphql`)

	http.ListenAndServe(":9090", nil)
}
