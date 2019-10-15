package store

// Store is an interface of key value store
type Store interface {
	Store(bucket, k string, v interface{})
	Load(bucket, k string) (interface{}, bool)
}
