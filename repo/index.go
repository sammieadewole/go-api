package repo

type ReadSource string

const (
	Postgres ReadSource = "postgres"
	Mongo    ReadSource = "mongo"
)

// Generic repository interface
type Repository[T any] interface {
	Create(entity *T) error
	Get() ([]*T, error)
	GetOne(id string) (*T, error)
	GetByEmail(email string) (*T, error)
	Update(id string, entity *T) error
	SoftDelete(id string) error
	HardDelete(id string) error
	AdminGet() ([]*T, error)
	AdminGetOne(id string) (*T, error)
	AdminGetByEmail(email string) (*T, error)
	AdminUpdate(id string, entity *T) error
}
