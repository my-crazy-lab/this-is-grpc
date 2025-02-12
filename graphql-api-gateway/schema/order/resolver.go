package schemaOrder

import (
	"log"

	"github.com/graphql-go/graphql"
	client "github.com/my-crazy-lab/this-is-grpc/proto-module/client"
	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
)

func AddToCartResolver(params graphql.ResolveParams) (interface{}, error) {
	user_id, _ := params.Args["user_id"].(int)
	product_id, _ := params.Args["product_id"].(int)
	quantity, _ := params.Args["quantity"].(int)

	ctx := params.Context

	res, err := addToCartGrpcHandler(ctx, client.AuthenticationService, &order.AddToCartRequest{
		UserId: int32(user_id), ProductId: int32(product_id), Quantity: int32(quantity),
	})
	if err != nil {
		log.Fatalf("AuthenticationClient.AddToCart(_) = _, %v: ", err)
	}

	return res.Item, nil
}

func ViewCartResolver(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	user_id, _ := params.Args["user_id"].(int)

	res, err := viewCartGrpcHandler(ctx, client.AuthenticationService, &order.ViewCartRequest{UserId: int32(user_id)})
	if err != nil {
		log.Fatalf("AuthenticationClient.ViewCart(_) = _, %v: ", err)
	}

	return res, nil
}

func PlaceOrderResolver(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	user_id, _ := params.Args["user_id"].(int)
	cart_id, _ := params.Args["cart_id"].(int)

	res, err := placeOrderGrpcHandler(ctx, client.AuthenticationService, &order.PlaceOrderRequest{UserId: int32(user_id), CartId: int32(cart_id)})
	if err != nil {
		log.Fatalf("AuthenticationClient.PlaceOrder(_) = _, %v: ", err)
	}

	return res, nil
}

func UpdateOrderStatusResolver(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	order_id, _ := params.Args["order_id"].(int)
	status, _ := params.Args["status"].(string)

	res, err := updateOrderStatusGrpcHandler(ctx, client.AuthenticationService, &order.UpdateOrderStatusRequest{Status: status, OrderId: int32(order_id)})
	if err != nil {
		log.Fatalf("AuthenticationClient.UpdateOrderStatus(_) = _, %v: ", err)
	}

	return res, nil
}

func CancelOrderResolver(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	order_id, _ := params.Args["order_id"].(int)

	res, err := cancelOrderGrpcHandler(ctx, client.AuthenticationService, &order.CancelOrderRequest{OrderId: int32(order_id)})
	if err != nil {
		log.Fatalf("AuthenticationClient.CancelOrder(_) = _, %v: ", err)
	}

	return res, nil
}

func GetOrderResolver(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context
	order_id, _ := params.Args["order_id"].(int)

	res, err := getOrderGrpcHandler(ctx, client.AuthenticationService, &order.GetOrderRequest{OrderId: int32(order_id)})
	if err != nil {
		log.Fatalf("AuthenticationClient.GetOrder(_) = _, %v: ", err)
	}

	return res, nil
}

func CreateShippingResolver(params graphql.ResolveParams) (interface{}, error) {
	ctx := params.Context

	order_id, _ := params.Args["order_id"].(int)
	shippingAddress, _ := params.Args["address"].(map[string]interface{})
	address := &order.ShippingAddress{
		Id:      int32(shippingAddress["id"].(int)),
		UserId:  int32(shippingAddress["user_id"].(int)),
		Address: shippingAddress["address"].(string),
		City:    shippingAddress["city"].(string),
		State:   shippingAddress["state"].(string),
		Country: shippingAddress["country"].(string),
		ZipCode: shippingAddress["zip_code"].(string),
	}

	res, err := createShippingGrpcHandler(ctx, client.AuthenticationService, &order.CreateShippingRequest{
		OrderId: int32(order_id),
		Address: address,
	})
	if err != nil {
		log.Fatalf("AuthenticationClient.CreateShipping(_) = _, %v: ", err)
	}

	return res, nil
}
