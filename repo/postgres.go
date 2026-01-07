package repo

import (
	"go-api/db"
)

// Postgres's Generic repo struct
type PostgresRepo[T interface{}] struct {
	tableName string
}

// Creates a new postgres repo
func NewPostgresRepo[T any](tableName string) *PostgresRepo[T] {
	return &PostgresRepo[T]{tableName: tableName}
}

// Creates a new entity
func (r *PostgresRepo[T]) Create(entity *T) error {
	return db.DB.Create(entity).Error
}

// Gets all entities
func (r *PostgresRepo[T]) Get() ([]*T, error) {
	var entities []*T
	return entities, db.DB.Where("deleted_at IS NULL").Find(&entities).Error
}

func (r *PostgresRepo[T]) GetByID(id string) (*T, error) {
	var entity T
	return &entity, db.DB.Where("id = ?", id).First(&entity).Error
}

func (r *PostgresRepo[T]) GetByEmail(email string) (*T, error) {
	var entity T
	return &entity, db.DB.Where("email = ?", email).First(&entity).Error
}
