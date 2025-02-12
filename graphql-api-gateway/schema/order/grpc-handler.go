package schemaOrder

import (
	"context"
	"errors"
	"log"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/auth"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
	"github.com/my-crazy-lab/this-is-grpc/shared/constants"
	"google.golang.org/grpc/metadata"
)

func createShippingGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.CreateShippingRequest) (*order.CreateShippingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.CreateShipping(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.CreateShipping(_) = _, %v: ", err)
	}

	return resp, nil
}

func cancelOrderGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.CancelOrder(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.CancelOrder(_) = _, %v: ", err)
	}

	return resp, nil
}

func updateOrderStatusGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.UpdateOrderStatus(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.UpdateOrderStatus(_) = _, %v: ", err)
	}

	return resp, nil
}

func getOrderGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.GetOrderRequest) (*order.OrderItem, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.GetOrder(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.GetOrder(_) = _, %v: ", err)
	}

	return resp, nil
}

func placeOrderGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.PlaceOrder(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.PlaceOrder(_) = _, %v: ", err)
	}

	return resp, nil
}

func viewCartGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.ViewCartRequest) (*order.ViewCartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.ViewCart(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.ViewCart(_) = _, %v: ", err)
	}

	return resp, nil
}

func addToCartGrpcHandler(ctx context.Context, client auth.AuthClient, params *order.AddToCartRequest) (*order.AddToCartResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.TIMEOUT)
	defer cancel()

	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := client.AddToCart(ctx, params)
	if err != nil {
		log.Fatalf("AuthenticationClient.AddToCart(_) = _, %v: ", err)
	}

	return resp, nil
}
