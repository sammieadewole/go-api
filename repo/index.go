package repo

type ReadSource string

const (
	Postgres ReadSource = "postgres"
	Mongo    ReadSource = "mongo"
)

// Generic repository interface
type Repository[T interface{}] interface {
	Create(entity *T) error
	Get() ([]*T, error)
	GetByID(id string) (*T, error)
	GetByEmail(email string) (*T, error)
}
