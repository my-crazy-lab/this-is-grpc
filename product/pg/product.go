package pg

import (
	"context"
)

type InsertCategoriesParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func InsertCategories(name, description string) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2)"
	_, err := DBPool.Exec(context.Background(), query, name, description)

	return err
}
