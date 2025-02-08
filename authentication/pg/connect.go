package pg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DBPool *pgxpool.Pool // Global variable for connection pool

// InitDB initializes the PostgreSQL connection pool
func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	fmt.Println("✅ Database pool initialized!")
	DBPool = dbpool
}

// CloseDB cleans up the database connection pool
func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		fmt.Println("🚀 Database connection closed")
	}
}
