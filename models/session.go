package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         string     `json:"id" bson:"_id" gorm:"primaryKey"`
	CustomerID string     `json:"customer_id" bson:"customer_id" gorm:"index"`
	Email      string     `json:"email" bson:"email" gorm:"index"`
	Token      string     `json:"token" bson:"token" gorm:"unique;index"`
	TTL        time.Time  `json:"ttl" bson:"ttl"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty" gorm:"index"`
}

func NewSession(customerID, email, token string) *Session {
	return &Session{
		ID:         uuid.NewString(),
		CustomerID: customerID,
		Email:      email,
		Token:      token,
		TTL:        time.Now().Add(24 * 7 * time.Hour),
	}
}

// Implements soft delete function
func (session *Session) SetDeleted(deletedAt *time.Time) {
	session.DeletedAt = deletedAt
}

func (session *Session) GetID() string {
	return session.ID
}

func (session *Session) SetID() {
	session.ID = uuid.NewString()
}
