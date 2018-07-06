package database

type Databaser interface {
	Initialize(dbpath string) error
	Finalize() error
	Save(id string, value []byte) error
	Load(id string) ([]byte, error)
	ForEach(func(k, v []byte) error) error
}
