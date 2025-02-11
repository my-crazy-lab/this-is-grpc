package client

import (
	"flag"
	"log"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/utils"

	"google.golang.org/grpc"
)

// avoid services interface establish multi times
var OrderService order.OrderClient
var OrderClientConnection *grpc.ClientConn

var addrOrderService = flag.String("addrOrderService", "localhost:50053", "the address to connect to")

func NewOrderClient() {
	if OrderService != nil {
		// avoid interface leaking
		OrderService = nil
	}

	// Set up the credentials for the connection.
	opts := utils.GetOptsClient()
	conn, err := grpc.NewClient(*addrOrderService, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	OrderClientConnection = conn
	OrderService = order.NewOrderClient(conn)
}
