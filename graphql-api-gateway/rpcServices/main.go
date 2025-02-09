package rpcServices

import (
	"flag"
	"log"

	authPb "github.com/my-crazy-lab/this-is-grpc/graph-api-gateway/proto/account"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/examples/data"
)

// avoid services interface establish multi times
var AuthenticationService authPb.AccountClient
var ClientConnection *grpc.ClientConn

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func NewAuthenticationService() {
	if AuthenticationService != nil {
		return
	}

	// Set up the credentials for the connection.
	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchToken())}
	creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// oauth.TokenSource requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	ClientConnection = conn
	AuthenticationService = authPb.NewAccountClient(conn)
}

// fetchToken simulates a token lookup and omits the details of proper token
// acquisition. For examples of how to acquire an OAuth2 token, see:
// https://godoc.org/golang.org/x/oauth2
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
