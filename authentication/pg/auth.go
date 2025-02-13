package pg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/my-crazy-lab/this-is-grpc/proto-module/proto/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx"
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
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Convert secret to []byte
	secretKey := []byte(getJwtSecret())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func VerifyJWT(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(getJwtSecret()), nil // Use the same secret key for verification
	})

	if err != nil {
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, fmt.Errorf("invalid token claims")
}

func isErrNoRows(err error) bool {
	return errors.Is(errors.Unwrap(err), pgx.ErrNoRows) || err.Error() == "no rows in result set"
}

func GetUserByPhone(phone string) (*User, error) {
	query := "SELECT id, phone, password FROM users WHERE phone = $1"
	row := DBPool.QueryRow(context.Background(), query, phone)

	var user User
	err := row.Scan(&user.ID, &user.PhoneNumber, &user.Password)
	if err != nil {
		if isErrNoRows(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

func GetUserById(id int32) (*user.User, error) {
	query := "SELECT id, phone FROM users WHERE id = $1"
	row := DBPool.QueryRow(context.Background(), query, id)
	var user user.User
	err := row.Scan(&user.Id, &user.PhoneNumber)
	if err != nil {
		if isErrNoRows(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func InsertNewUser(phone, password string) error {
	passHashed, err := hashPassword(password)

	if err != nil {
		return err
	}

	query := "INSERT INTO users (phone, password) VALUES ($1, $2)"
	_, err = DBPool.Exec(context.Background(), query, phone, passHashed)

	return err
}

func GetUsers() ([]User, error) {
	query := "SELECT id, phone FROM users"
	rows, err := DBPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// var users []User
	users := make([]User, 0) // ✅ Always initialize slice

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.PhoneNumber); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
