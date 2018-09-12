package database

// Databaser is an interface to treat databases
type Databaser interface {
	Initialize(dbpath string) error
	Finalize() error
	Save(bucket, id string, value []byte) error
	Load(bucket, id string) ([]byte, error)
	ForEach(bucket string, f func(k, v []byte) error) error
}
