package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PlaceOrder(req *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	ctx := context.Background()
	tx, err := DBPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if cart exists and is active
	var cartStatus string
	err = tx.QueryRow(ctx, `
		SELECT status FROM carts WHERE id = $1 AND user_id = $2
	`, req.CartId, req.UserId).Scan(&cartStatus)
	if err != nil {
		return nil, fmt.Errorf("cart not found or not active: %w", err)
	}
	if cartStatus != "active" {
		return nil, fmt.Errorf("cart is not active")
	}

	// Get cart items and check inventory
	rows, err := tx.Query(ctx, `
		SELECT ci.product_id, ci.quantity, i.quantity
		FROM cart_items ci
		JOIN inventory i ON ci.product_id = i.product_id
		WHERE ci.cart_id = $1
	`, req.CartId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	}
	defer rows.Close()

	var items []struct {
		ProductID       int32
		QuantityNeeded  int32
		QuantityInStock int32
	}

	for rows.Next() {
		var item struct {
			ProductID       int32
			QuantityNeeded  int32
			QuantityInStock int32
		}
		if err := rows.Scan(&item.ProductID, &item.QuantityNeeded, &item.QuantityInStock); err != nil {
			return nil, fmt.Errorf("error scanning cart item: %w", err)
		}
		if item.QuantityNeeded > item.QuantityInStock {
			return nil, fmt.Errorf("insufficient inventory for product_id: %d", item.ProductID)
		}
		items = append(items, item)
	}

	// Get current timestamp
	now := time.Now()

	// Create order
	var orderID int32
	var createdAt time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, cart_id, total, status, created_at, updated_at)
		VALUES ($1, $2, 0, 'pending', $3, $3) RETURNING id, created_at
	`, req.UserId, req.CartId, now).Scan(&orderID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Insert order items and calculate total
	var total int32
	for _, item := range items {
		var price int32
		err = tx.QueryRow(ctx, `SELECT price FROM products WHERE id = $1`, item.ProductID).Scan(&price)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product price: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, price, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, orderID, item.ProductID, item.QuantityNeeded, price, now)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item: %w", err)
		}

		total += item.QuantityNeeded * price
	}

	// Update order total
	_, err = tx.Exec(ctx, `UPDATE orders SET total = $1, updated_at = $2 WHERE id = $3`, total, now, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to update order total: %w", err)
	}

	// Reduce inventory
	for _, item := range items {
		_, err = tx.Exec(ctx, `
			UPDATE inventory SET quantity = quantity - $1 WHERE product_id = $2
		`, item.QuantityNeeded, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to update inventory: %w", err)
		}
	}

	// Mark cart as checked out
	_, err = tx.Exec(ctx, `UPDATE carts SET status = 'checkout', updated_at = $1 WHERE id = $2`, now, req.CartId)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return response
	return &order.PlaceOrderResponse{
		Order: &order.OrderItem{
			Id:        orderID,
			UserId:    req.UserId,
			CartId:    req.CartId,
			Total:     total,
			Status:    "pending",
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(now),
		},
	}, nil
}

func UpdateOrderStatus(req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	ctx := context.Background()

	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}
	if !validStatuses[req.Status] {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}

	// Start transaction
	tx, err := DBPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update order status
	now := time.Now()
	_, err = tx.Exec(ctx, `
		UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3
	`, req.Status, now, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Fetch updated order details
	var orderItem order.OrderItem
	err = tx.QueryRow(ctx, `
		SELECT id, user_id, cart_id, total, status, created_at, updated_at
		FROM orders WHERE id = $1
	`, req.OrderId).Scan(
		&orderItem.Id,
		&orderItem.UserId,
		&orderItem.CartId,
		&orderItem.Total,
		&orderItem.Status,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated order: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return updated order response
	return &order.UpdateOrderStatusResponse{
		Order: &order.OrderItem{
			Id:        orderItem.Id,
			UserId:    orderItem.UserId,
			CartId:    orderItem.CartId,
			Total:     orderItem.Total,
			Status:    orderItem.Status,
			CreatedAt: timestamppb.New(orderItem.CreatedAt.AsTime()),
			UpdatedAt: timestamppb.New(now),
		},
	}, nil
}

func CancelOrder(req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	ctx := context.Background()

	// Start transaction
	tx, err := DBPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Fetch order details and ensure it can be canceled
	var orderItem order.OrderItem
	err = tx.QueryRow(ctx, `
		SELECT id, user_id, cart_id, total, status, created_at, updated_at
		FROM orders WHERE id = $1
	`, req.OrderId).Scan(
		&orderItem.Id,
		&orderItem.UserId,
		&orderItem.CartId,
		&orderItem.Total,
		&orderItem.Status,
		&orderItem.CreatedAt,
		&orderItem.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Only allow cancellation if order is not already shipped or delivered
	if orderItem.Status == "shipped" || orderItem.Status == "delivered" {
		return nil, fmt.Errorf("order cannot be canceled as it is already %s", orderItem.Status)
	}

	// Update order status to "cancelled"
	now := time.Now()
	_, err = tx.Exec(ctx, `
		UPDATE orders SET status = 'cancelled', updated_at = $1 WHERE id = $2
	`, now, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Restore inventory if order was already deducted
	rows, err := tx.Query(ctx, `
		SELECT product_id, quantity FROM order_items WHERE order_id = $1
	`, req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var productID, quantity int
		if err := rows.Scan(&productID, &quantity); err != nil {
			return nil, fmt.Errorf("failed to read order item: %w", err)
		}

		_, err = tx.Exec(ctx, `
			UPDATE inventory SET quantity = quantity + $1 WHERE product_id = $2
		`, quantity, productID)
		if err != nil {
			return nil, fmt.Errorf("failed to restore inventory: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return updated order response
	return &order.CancelOrderResponse{
		Order: &order.OrderItem{
			Id:        orderItem.Id,
			UserId:    orderItem.UserId,
			CartId:    orderItem.CartId,
			Total:     orderItem.Total,
			Status:    "cancelled",
			CreatedAt: timestamppb.New(orderItem.CreatedAt.AsTime()),
			UpdatedAt: timestamppb.New(now),
		},
	}, nil
}

func GetOrder(req *order.GetOrderRequest) (*order.OrderItem, error) {
	ctx := context.Background()

	// Query the order details
	var orderItem order.OrderItem
	var createdAt, updatedAt time.Time

	err := DBPool.QueryRow(ctx, `
		SELECT id, user_id, cart_id, total, status, created_at, updated_at
		FROM orders WHERE id = $1
	`, req.OrderId).Scan(
		&orderItem.Id,
		&orderItem.UserId,
		&orderItem.CartId,
		&orderItem.Total,
		&orderItem.Status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Convert timestamps to protobuf format
	orderItem.CreatedAt = timestamppb.New(createdAt)
	orderItem.UpdatedAt = timestamppb.New(updatedAt)

	return &orderItem, nil
}

func CreateShipping(req *order.CreateShippingRequest) (*order.CreateShippingResponse, error) {
	ctx := context.Background()
	tx, err := DBPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert shipping address
	const SqlInsertShippingAddress = `
		INSERT INTO shipping_addresses 
		(user_id, address, city, state, country, zip_code, created_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) 
		RETURNING
		id, created_at
	`
	var addressID int32
	var createdAt, updatedAt time.Time

	err = tx.QueryRow(ctx, SqlInsertShippingAddress,
		req.Address.UserId, req.Address.Address, req.Address.City,
		req.Address.State, req.Address.Country, req.Address.ZipCode).
		Scan(&addressID, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert shipping address: %w", err)
	}

	// Insert shipping record
	const SqlInsertShippingRecord = `
		INSERT INTO shippings (order_id, shipping_address_id, status, created_at, updated_at) 
		VALUES ($1, $2, 'pending', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
		RETURNING id, created_at, updated_at
	`
	var shippingID int32
	err = tx.QueryRow(ctx, SqlInsertShippingRecord,
		req.OrderId, addressID).
		Scan(&shippingID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create shipping: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Construct response
	response := &order.CreateShippingResponse{
		Shipping: &order.Shipping{
			Id:      shippingID,
			OrderId: req.OrderId,
			Address: &order.ShippingAddress{
				Id:        addressID,
				UserId:    req.Address.UserId,
				Address:   req.Address.Address,
				City:      req.Address.City,
				State:     req.Address.State,
				Country:   req.Address.Country,
				ZipCode:   req.Address.ZipCode,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
			Status:    "pending",
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
	}

	return response, nil
}

func ViewCart(req *order.ViewCartRequest) (*order.ViewCartResponse, error) {
	ctx := context.Background()

	// Get the active cart for the user
	const SqlGetActiveCartByUserId = `
		SELECT id, user_id, status, created_at, updated_at 
		FROM carts 
		WHERE user_id = $1 AND status = 'active' 
		LIMIT 1
	`
	var cart order.Cart
	err := DBPool.QueryRow(ctx, SqlGetActiveCartByUserId, req.UserId).Scan(
		&cart.Id, &cart.UserId, &cart.Status,
		&cart.CreatedAt, &cart.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}

	// Get the cart items
	const SqlGetCartItems = `
		SELECT id, cart_id, product_id, quantity, created_at 
		FROM cart_items 
		WHERE cart_id = $1
	`
	rows, err := DBPool.Query(ctx, SqlGetCartItems, cart.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	}
	defer rows.Close()

	var items []*order.CartItem
	for rows.Next() {
		var item order.CartItem
		err := rows.Scan(&item.Id, &item.CartId, &item.ProductId, &item.Quantity, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart item: %w", err)
		}
		items = append(items, &item)
	}

	// Return the cart and its items
	return &order.ViewCartResponse{
		Cart:  &cart,
		Items: items,
	}, nil
}

func AddToCart(req *order.AddToCartRequest) (*order.AddToCartResponse, error) {
	productId := req.ProductId
	userId := req.UserId
	quantity := req.Quantity

	if DBPool == nil {
		return nil, status.Errorf(codes.Internal, "database pool is nil")
	}

	tx, err := DBPool.Begin(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Step 1: Check inventory
	var currentQuantity int32
	err = tx.QueryRow(context.Background(), "SELECT quantity FROM inventory WHERE product_id = $1", productId).Scan(&currentQuantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to fetch inventory: %v", err)
	}

	if quantity > currentQuantity {
		return nil, status.Errorf(codes.FailedPrecondition, "Not enough stock available")
	}

	// Step 2: Get or create active cart
	var cartID int32
	err = tx.QueryRow(context.Background(), `  
        SELECT id FROM carts   
        WHERE user_id = $1 AND status = 'active';  
    `, userId).Scan(&cartID)

	if err != nil {
		// If no active cart exists, create a new one
		_, err = tx.Exec(context.Background(), `  
            INSERT INTO carts (user_id, status, created_at)  
            VALUES ($1, 'active', NOW())  
            RETURNING id;  
        `, userId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to create new cart: %v", err)
		}
		if err := tx.QueryRow(context.Background(), "SELECT currval('carts_id_seq'::regclass)").Scan(&cartID); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to get new cart ID: %v", err)
		}
	}

	// Step 3: Add or update cart item with proper ON CONFLICT
	_, err = tx.Exec(context.Background(), `  
        INSERT INTO cart_items (cart_id, product_id, quantity, created_at)  
        VALUES ($1, $2, $3, NOW())  
        ON CONFLICT (cart_id, product_id)   
        DO UPDATE SET quantity = cart_items.quantity + $3;  
    `, cartID, productId, quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to add/update cart item: %v", err)
	}

	// Step 4: Get the item ID
	var itemID int32
	err = tx.QueryRow(context.Background(), `  
        SELECT id FROM cart_items  
        WHERE cart_id = $1 AND product_id = $2;  
    `, cartID, productId).Scan(&itemID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get cart item ID: %v", err)
	}

	// Step 5: Commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	// Prepare the response
	cartItem := &order.CartItem{
		Id:        itemID,
		CartId:    cartID,
		ProductId: productId,
		Quantity:  quantity,
		CreatedAt: timestamppb.New(time.Now()),
	}

	return &order.AddToCartResponse{Item: cartItem}, nil
}
