package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	User  Role = "user"
	Admin Role = "admin"
)

type Customer struct {
	ID             string     `json:"id" bson:"_id" gorm:"primaryKey"`
	Name           string     `json:"name" bson:"name" gorm:"not null"`
	Email          string     `json:"email" bson:"email" gorm:"unique;not null"`
	Phone          string     `json:"phone" bson:"phone"`
	HashedPassword string     `json:"-" bson:"password" gorm:"not null"`
	Role           Role       `json:"role" bson:"role" gorm:"default:'user'"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty" gorm:"index"`
}

func (customer *Customer) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(customer.HashedPassword), []byte(password)) == nil
}

// Implements soft delete function
func (customer *Customer) SetDeleted(deletedAt *time.Time) {
	customer.DeletedAt = deletedAt
}

func (customer *Customer) GetID() string {
	return customer.ID
}

func (customer *Customer) SetID() {
	customer.ID = uuid.NewString()
}

func NewCustomer(name, email, phone string) (*Customer, error) {

	return &Customer{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

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
