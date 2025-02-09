package schema

import "github.com/graphql-go/graphql"

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"todo":     todoQuery["todo"],
		"lastTodo": todoQuery["lastTodo"],
		"todoList": todoQuery["todoList"],
		"getUsers": authQuery["getUsers"],
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"createTodo": todoMutation["createTodo"],
		"updateTodo": todoMutation["updateTodo"],
		"signIn":     authMutation["signIn"],
		"signUp":     authMutation["signUp"],
	},
})

var RootSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
