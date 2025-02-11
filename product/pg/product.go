package pg

import (
	"encoding/json"
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
			CreatedAt:   timestamppb.New(createdAt),
			UpdatedAt:   timestamppb.New(updatedAt),
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

const SqlFetchProductWithPagination = `
	SELECT 
		p.id, p.name, p.description, p.price, p.created_at, p.updated_at, 
		COALESCE(i.quantity, 0) AS quantity,
		COALESCE(json_agg(json_build_object(
			'id', c.id, 
			'name', c.name, 
			'description', c.description, 
			'created_at', c.created_at::TEXT, 
			'updated_at', c.updated_at::TEXT)) 
			FILTER (WHERE c.id IS NOT NULL), '[]') AS categories,
		COUNT(*) OVER() AS total_count
	FROM products p
	LEFT JOIN product_categories pc ON p.id = pc.product_id
	LEFT JOIN categories c ON pc.category_id = c.id
	LEFT JOIN inventory i ON p.id = i.product_id
	GROUP BY p.id, i.quantity
	ORDER BY p.created_at DESC
	LIMIT $1 OFFSET $2;
`
const LAYOUT = "2006-01-02 15:04:05.999999-07"

func GetProducts(pageSize, pageIndex int32) ([]*productPb.ProductItem, int32, error) {
	offset := (pageIndex - 1) * pageSize
	rows, err := DBPool.Query(context.Background(), SqlFetchProductWithPagination, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*productPb.ProductItem
	var total int32

	for rows.Next() {
		var (
			createdAt      time.Time
			updatedAt      time.Time
			categoriesJSON string
		)

		var product productPb.ProductItem

		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &createdAt, &updatedAt, &product.Quantity, &categoriesJSON, &total)
		if err != nil {
			return nil, 0, err
		}

		product.CreatedAt = timestamppb.New(createdAt)
		product.UpdatedAt = timestamppb.New(updatedAt)

		// Convert JSON categories to Go struct
		var categoryList []struct {
			Id          int32  `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			CreatedAt   string `json:"created_at"`
			UpdatedAt   string `json:"updated_at"`
		}
		if err := json.Unmarshal([]byte(categoriesJSON), &categoryList); err != nil {
			return nil, 0, err
		}
		// Convert to protobuf format
		var protoCategories []*productPb.Category
		for _, c := range categoryList {
			// Convert into time.Time
			createdAt, _ := time.Parse(LAYOUT, c.CreatedAt)
			updatedAt, _ := time.Parse(LAYOUT, c.UpdatedAt)

			protoCategories = append(protoCategories, &productPb.Category{
				Id:          c.Id,
				Name:        c.Name,
				Description: c.Description,
				CreatedAt:   timestamppb.New(createdAt),
				UpdatedAt:   timestamppb.New(updatedAt),
			})
		}

		product.Categories = protoCategories
		products = append(products, &product)
	}

	return products, total, nil
}
