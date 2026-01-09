package middleware

import (
	"errors"
	"fmt"
	"go-api/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Create Token
func GenerateToken(userID, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := models.Claims{
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

	claims := models.Claims{
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

// // Validate Token
// func ValidateToken(authToken string) (*models.Claims, error) {
// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		return nil, errors.New("JWT_SECRET not set")
// 	}

// 	token, err := jwt.ParseWithClaims(authToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		// Verify signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid signing method")
// 		}
// 		return []byte(secret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if !token.Valid {
// 		return nil, errors.New("invalid token")
// 	}

// 	claims, ok := token.Claims.(*models.Claims)
// 	if !ok {
// 		return nil, errors.New("invalid token claims")
// 	}

// 	return claims, nil
// }

func VerifyIDToken(ctx *gin.Context, tokenString string) (*models.Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		fmt.Println("JWT parse error:", err)
		fmt.Printf("Token valid: %v", token.Valid)
		fmt.Printf("Error: %v", err)
		fmt.Printf("Invalid or expired token: %v", token.Claims)
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Invalid claims format")
		return nil, errors.New("invalid claims format")
	}

	email, _ := claims["email"].(string)

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		fmt.Println("Missing exp claim")
		return nil, errors.New("missing exp claim")
	}
	expirationTime := time.Unix(int64(expFloat), 0)
	if time.Now().After(expirationTime) {
		fmt.Println("Token expired at:", expirationTime)
		return nil, errors.New("token has expired")
	}

	return &models.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}, nil
}
