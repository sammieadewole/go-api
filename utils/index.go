package utils

import (
	"errors"
	"go-api/models"
	"go-api/repo"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Safe customer info to be sent to clients (no ID or IsDeleted)
type CustomerPublic struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Claims for jwt
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// StripSensitiveData removes ID/IsDeleted and masks email/phone
func StripSensitiveData(customer *models.Customer) CustomerPublic {
	return CustomerPublic{
		Name:      customer.Name,
		Email:     maskString(customer.Email),
		Phone:     maskString(customer.Phone),
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

// StripSensitiveDataList processes multiple customers
func StripSensitiveDataList(customers []*models.Customer) []CustomerPublic {
	result := make([]CustomerPublic, len(customers))
	for i, customer := range customers {
		result[i] = StripSensitiveData(customer)
	}
	return result
}

// maskString safely masks the middle of a string
func maskString(s string) string {
	if len(s) <= 6 {
		return "***"
	}
	return s[:3] + "..." + s[len(s)-3:]
}

func GetReadSource(c *gin.Context) repo.ReadSource {
	storage := c.Query("storage")
	if storage == "sql" {
		return repo.Postgres
	}
	return repo.Mongo
}

// Create Token
func GenerateToken(userID, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Function to create a token with expiry data
func GenerateTokenWithExpiry(userID, email string, expiry time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Validate Token
func ValidateToken(authToken string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(authToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// Create a password hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Validate password with hash
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
