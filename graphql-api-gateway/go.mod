module github.com/my-crazy-lab/this-is-grpc/graphql-api-gateway

go 1.23.6

replace github.com/my-crazy-lab/this-is-grpc/proto-module => ../proto-module

require (
	github.com/graphql-go/graphql v0.8.1
	github.com/my-crazy-lab/this-is-grpc/proto-module v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.26.0
	google.golang.org/grpc v1.70.0
	google.golang.org/grpc/examples v0.0.0-20250207054034-9afb49d37828
	google.golang.org/protobuf v1.36.5
)

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)
