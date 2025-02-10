package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/rpcServices"
	"github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway/schema"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operationName"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	rpcServices.NewAuthenticationService()
	defer rpcServices.AuthClientConnection.Close()

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

		result := graphql.Do(graphql.Params{
			Context:        req.Context(),
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

	fmt.Println(`Get users:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "{ getUsers { id phoneNumber password } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Sign in:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "mutation { signIn(phoneNumber:\"123456\",password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Sign up:
curl \
-X POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer _token_" \
--data '{ "query": "mutation { signUp(phoneNumber:\"123456\", password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql`)

	http.ListenAndServe(":9090", nil)
}
