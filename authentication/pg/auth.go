package pg

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func getJwtSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("JWT_SECRET")
}

type User struct {
	ID int `json:"id"`
	LoginParams
}

type LoginParams struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

// JWT Claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Hash password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT Token
func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtSecret())
}

func GetUserByPhone(phone string) (*User, error) {
	query := "SELECT id, phone, password FROM users WHERE phone = $1"
	row := DBPool.QueryRow(context.Background(), query, phone)
	fmt.Printf("phone: ", phone)
	var user User
	if err := row.Scan(&user.ID, &user.PhoneNumber, &user.Password); err != nil {
		return nil, err
	}
	fmt.Printf(user.Password)
	return &user, nil
}
