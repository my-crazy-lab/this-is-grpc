package pg

import (
	"fmt"
	"time"

	productPb "github.com/my-crazy-lab/this-is-grpc/proto-module/proto/product"
	"google.golang.org/protobuf/types/known/timestamppb"

	"context"
)

type InsertCategoriesParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func CreateProduct(name, description string, price float64, categories *[]int32, quantity int32) (int32, error) {
	ctx := context.Background()
	tx, err := DBPool.Begin(ctx) // Start transaction
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx) // Rollback on error
		} else {
			tx.Commit(ctx) // Commit if successful
		}
	}()

	// Insert product
	var productID int32
	insertProductQuery := `
		INSERT INTO products (name, description, price) 
		VALUES ($1, $2, $3) 
		RETURNING id`
	err = tx.QueryRow(ctx, insertProductQuery, name, description, price).Scan(&productID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert product: %w", err)
	}

	// Insert inventory
	insertInventoryQuery := `
		INSERT INTO inventory (product_id, quantity) 
		VALUES ($1, $2)`
	_, err = tx.Exec(ctx, insertInventoryQuery, productID, quantity)
	if err != nil {
		return 0, fmt.Errorf("failed to insert inventory: %w", err)
	}

	// Insert product categories (if any)
	if len(*categories) > 0 {
		insertCategoryQuery := `
			INSERT INTO product_categories (product_id, category_id)
			VALUES ($1, $2)`
		for _, categoryID := range *categories {
			_, err = tx.Exec(ctx, insertCategoryQuery, productID, categoryID)
			if err != nil {
				return 0, fmt.Errorf("failed to insert product-category mapping: %w", err)
			}
		}
	}

	// If we reach here, everything is successful, commit transaction
	return productID, nil
}

func InsertCategories(name, description string) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2)"
	_, err := DBPool.Exec(context.Background(), query, name, description)

	return err
}

func GetCategories() ([]*productPb.Category, error) {
	query := "SELECT id, name, description, created_at, updated_at FROM categories"
	rows, err := DBPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*productPb.Category, 0)

	for rows.Next() {
		var (
			id          int32
			name        string
			description string
			createdAt   time.Time
			updatedAt   time.Time
		)

		// Scan into variables
		if err := rows.Scan(&id, &name, &description, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		// Convert time.Time to *timestamppb.Timestamp
		category := &productPb.Category{
			Id:          id,
			Name:        name,
			Description: description,
			CreatedAt:   timestamppb.New(createdAt), // ✅ Convert correctly
			UpdatedAt:   timestamppb.New(updatedAt), // ✅ Convert correctly
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
