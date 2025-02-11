package pg

import (
	"context"
	"time"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AddToCart(productId, userId, quantity int32) (*order.AddToCartResponse, error) {
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
