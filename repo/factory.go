package repo

import (
	"fmt"
	"time"
)

// RepoFactory Type for repositories
// Primary and secondary can be Postgres or Mongo in any other
type RepoFactory[T any] struct {
	primary   Repository[T]
	secondary Repository[T]
}

// Create a new RepoFactory
func NewRepoFactory[T any](primary, secondary Repository[T]) *RepoFactory[T] {
	return &RepoFactory[T]{
		primary:   primary,
		secondary: secondary,
	}
}

// Creates a new entity, writes to primary first, then secondary
// Rolls back if secondary fails
func (f *RepoFactory[T]) Create(entity *T) error {
	if err := f.primary.Create(entity); err != nil {
		return fmt.Errorf("primary create failed: %w", err)
	}

	if err := f.secondary.Create(entity); err != nil {
		f.primary.HardDelete(getEntityID(entity))
		return fmt.Errorf("secondary create failed, rolled back: %w", err)
	}

	return nil
}

// Updates an entity, writes to primary first, then secondary
// Rolls back if secondary fails
func (f *RepoFactory[T]) Update(id string, entity *T) error {
	original, err := f.primary.GetOne(id)
	if err != nil {
		return fmt.Errorf("get original failed: %w", err)
	}

	if err := f.primary.Update(id, entity); err != nil {
		return fmt.Errorf("primary update failed: %w", err)
	}

	if err := f.secondary.Update(id, entity); err != nil {
		f.primary.Update(id, original)
		return fmt.Errorf("secondary update failed, rolled back: %w", err)
	}

	return nil
}

// Soft deletes an entity, writes to primary first, then secondary
// Rolls back if secondary fails
func (f *RepoFactory[T]) SoftDelete(id string) error {
	if err := f.primary.SoftDelete(id); err != nil {
		return fmt.Errorf("primary soft delete failed: %w", err)
	}

	if err := f.secondary.SoftDelete(id); err != nil {
		if entity, getErr := f.primary.GetOne(id); getErr == nil {
			setEntityDeleted(entity, nil)
			f.primary.Update(id, entity)
		}
		return fmt.Errorf("secondary soft delete failed, rolled back: %w", err)
	}

	return nil
}

// Hard deletes an entity, writes to primary first, then secondary
// Rolls back if secondary fails
func (f *RepoFactory[T]) HardDelete(id string) error {
	original, err := f.primary.AdminGetOne(id)
	if err != nil {
		return fmt.Errorf("get original failed: %w", err)
	}

	if err := f.primary.HardDelete(id); err != nil {
		return fmt.Errorf("primary hard delete failed: %w", err)
	}

	if err := f.secondary.HardDelete(id); err != nil {
		f.primary.Create(original)
		return fmt.Errorf("secondary hard delete failed, rolled back: %w", err)
	}

	return nil
}

// Gets all entities from either primary or secondary
// Defaults to primary
func (f *RepoFactory[T]) Get(readFrom ReadSource) ([]*T, error) {
	if readFrom == Mongo {
		return f.secondary.Get()
	}
	return f.primary.Get()
}

// Gets one entity from either primary or secondary
// Defaults to primary
func (f *RepoFactory[T]) GetOne(id string, readFrom ReadSource) (*T, error) {
	if readFrom == Mongo {
		return f.secondary.GetOne(id)
	}
	return f.primary.GetOne(id)
}

// Gets one entity from either primary or secondary
// Defaults to primary
func (f *RepoFactory[T]) GetByEmail(email string, readFrom ReadSource) (*T, error) {
	if readFrom == Mongo {
		return f.secondary.GetByEmail(email)
	}
	return f.primary.GetByEmail(email)
}

// Admin functions

// Gets all entities from either primary or secondary whether soft deleted or not
// Defaults to primary
func (f *RepoFactory[T]) AdminGet(readFrom ReadSource) ([]*T, error) {
	if readFrom == Mongo {
		return f.secondary.AdminGet()
	}
	return f.primary.AdminGet()
}

// Gets one entity from either primary or secondary whether soft deleted or not
// Defaults to primary
func (f *RepoFactory[T]) AdminGetOne(id string, readFrom ReadSource) (*T, error) {
	if readFrom == Mongo {
		return f.secondary.AdminGetOne(id)
	}
	return f.primary.AdminGetOne(id)
}

// Gets one entity from either primary or secondary whether soft deleted or not
// Defaults to primary
func (f *RepoFactory[T]) AdminGetByEmail(email string, readFrom ReadSource) (*T, error) {
	if readFrom == Mongo {
		return f.secondary.AdminGetByEmail(email)
	}
	return f.primary.AdminGetByEmail(email)
}

// Updates an entity whether soft deleted or not
// Writes to primary first, then secondary
// Rolls back if secondary fails
func (f *RepoFactory[T]) AdminUpdate(id string, entity *T) error {
	original, err := f.primary.AdminGetOne(id)
	if err != nil {
		return fmt.Errorf("get original failed: %w", err)
	}

	if err := f.primary.AdminUpdate(id, entity); err != nil {
		return fmt.Errorf("primary admin update failed: %w", err)
	}

	if err := f.secondary.AdminUpdate(id, entity); err != nil {
		f.primary.AdminUpdate(id, original)
		return fmt.Errorf("secondary admin update failed, rolled back: %w", err)
	}

	return nil
}

// Helpers

// Gets the ID of an entity
func getEntityID(entity interface{}) string {
	type IDGetter interface {
		GetID() string
	}
	if e, ok := entity.(IDGetter); ok {
		return e.GetID()
	}
	return ""
}

// Sets the ID of an entity
func setEntityID(entity interface{}, id string) {
	type IDSetter interface {
		SetID(string)
	}
	if e, ok := entity.(IDSetter); ok {
		e.SetID(id)
	}
}

// Soft deletes an entity
func setEntityDeleted(entity interface{}, deletedAt *time.Time) {
	type DeletedSetter interface {
		SetDeleted(*time.Time)
	}
	if e, ok := entity.(DeletedSetter); ok {
		e.SetDeleted(deletedAt)
	}
}
