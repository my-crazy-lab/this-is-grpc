package client

import (
	"flag"
	"log"

	authPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/utils"

	"google.golang.org/grpc"
)

// avoid services interface establish multi times
var AuthenticationService authPb.AuthClient
var AuthClientConnection *grpc.ClientConn

var addrAuthService = flag.String("addrAuthService", "localhost:50051", "the address to connect to")

func NewAuthenticationClient() {
	if AuthenticationService != nil {
		// avoid interface leaking
		AuthenticationService = nil
	}

	// Set up the credentials for the connection.
	opts := utils.GetOptsClient()
	conn, err := grpc.NewClient(*addrAuthService, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	AuthClientConnection = conn
	AuthenticationService = authPb.NewAuthClient(conn)
}
