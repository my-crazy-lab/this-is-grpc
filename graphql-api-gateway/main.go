package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/rpcServices"
	"github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/schema"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operationName"`
	Variables map[string]interface{} `json:"variables"`
}

func main() {
	rpcServices.NewAuthenticationService()
	defer rpcServices.ClientConnection.Close()

	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p postData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			w.WriteHeader(400)
			return
		}
		result := graphql.Do(graphql.Params{
			Context:        req.Context(),
			Schema:         schema.RootSchema,
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
--data '{ "query": "{ getUsers { id phoneNumber password } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Create new todo:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "mutation { createTodo(text:\"My New todo\") { id text done } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Update todo:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "mutation { updateTodo(id:\"a\", done: true) { id text done } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Load todo list:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "{ todoList { id text done } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Sign in:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "mutation { signIn(phoneNumber:\"123456\",password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql`)

	fmt.Println("")

	fmt.Println(`Sign up:
curl \
-X POST \
-H "Content-Type: application/json" \
--data '{ "query": "mutation { signUp(phoneNumber:\"123456\", password:\"hihihi\") { token } }" }' \
http://localhost:9090/graphql`)

	http.ListenAndServe(":9090", nil)
}
