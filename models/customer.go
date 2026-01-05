package models

import (
	"time"

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

func NewCustomer(name, email, phone, password string) (*Customer, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return nil, err
	}

	return &Customer{
		ID:             uuid.NewString(),
		Name:           name,
		Email:          email,
		Phone:          phone,
		HashedPassword: string(hashedPassword),
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
