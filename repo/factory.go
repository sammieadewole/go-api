package repo

// RepoFactory Type for repositories
// Primary and secondary can be Postgres or Mongo in any other
type RepoFactory[T any] struct {
	postgres Repository[T]
	mongo    Repository[T]
}

// Create a new RepoFactory
func NewRepoFactory[T any](postgres, mongo Repository[T]) *RepoFactory[T] {
	return &RepoFactory[T]{
		postgres: postgres,
		mongo:    mongo,
	}
}

// Creates a new entity, writes to primary first, then secondary
// Rolls back if secondary fails
func (f *RepoFactory[T]) Create(entity *T, readFrom ReadSource) error {
	if readFrom == Mongo {
		return f.mongo.Create(entity)
	}
	return f.postgres.Create(entity)
}

// Gets all entities from either primary or secondary
// Defaults to primary
func (f *RepoFactory[T]) Get(readFrom ReadSource) ([]*T, error) {
	if readFrom == Mongo {
		return f.mongo.Get()
	}
	return f.postgres.Get()
}

func (f *RepoFactory[T]) GetByID(id string, readFrom ReadSource) (*T, error) {
	if readFrom == Mongo {
		return f.mongo.GetByID(id)
	}
	return f.postgres.GetByID(id)
}

func (f *RepoFactory[T]) GetByEmail(email string, readFrom ReadSource) (*T, error) {
	if readFrom == Mongo {
		return f.mongo.GetByEmail(email)
	}
	return f.postgres.GetByEmail(email)
}
