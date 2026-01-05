package repo

import (
	"go-api/db"
	"time"
)

// Postgres's Generic repo struct
type PostgresRepo[T any] struct {
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

// Gets one entity by ID
func (r *PostgresRepo[T]) GetOne(id string) (*T, error) {
	var entity T
	if err := db.DB.Where("deleted_at IS NULL").First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Gets one entity by email
func (r *PostgresRepo[T]) GetByEmail(email string) (*T, error) {
	var entity T
	if err := db.DB.Where("deleted_at IS NULL").First(&entity, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Updates an entity
func (r *PostgresRepo[T]) Update(id string, entity *T) error {
	return db.DB.Where("deleted_at IS NULL").Save(entity).Error
}

// Soft deletes an entity
func (r *PostgresRepo[T]) SoftDelete(id string) error {
	var entity T

	if err := db.DB.First(&entity, "id = ?", id).Error; err != nil {
		return err
	}

	now := time.Now()
	setEntityDeleted(&entity, &now)
	return db.DB.Save(&entity).Error
}

// Hard deletes an entity
func (r *PostgresRepo[T]) HardDelete(id string) error {
	var entity T
	return db.DB.Unscoped().Delete(&entity, "id = ?", id).Error
}

// Admin functions

// Gets all entities whether soft deleted or not
func (r *PostgresRepo[T]) AdminGet() ([]*T, error) {
	var entities []*T
	return entities, db.DB.Unscoped().Find(&entities).Error
}

// Gets one entity by id whether soft deleted or not
func (r *PostgresRepo[T]) AdminGetOne(id string) (*T, error) {
	var entity T
	return &entity, db.DB.Unscoped().First(&entity, "id = ?", id).Error
}

// Gets one entity by email whether soft deleted or not
func (r *PostgresRepo[T]) AdminGetByEmail(email string) (*T, error) {
	var entity T
	return &entity, db.DB.Unscoped().First(&entity, "email = ?", email).Error
}

// Updates an entity whether soft deleted or not
func (r *PostgresRepo[T]) AdminUpdate(id string, entity *T) error {
	return db.DB.Save(entity).Error
}
