package client

import (
	"flag"
	"log"

	productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/utils"

	"google.golang.org/grpc"
)

// avoid services interface establish multi times
var ProductService productPb.ProductClient
var ProductClientConnection *grpc.ClientConn

var addrProductService = flag.String("addrProductService", "localhost:50052", "the address to connect to")

func NewProductClient() {
	if ProductService != nil {
		// avoid interface leaking
		ProductService = nil
	}

	// Set up the credentials for the connection.
	opts := utils.GetOptsClient()
	conn, err := grpc.NewClient(*addrProductService, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	ProductClientConnection = conn
	ProductService = productPb.NewProductClient(conn)
}
