package kvstore

// KVStore is an interface of key value store
type KVStore interface {
	Store(bucket, k string, v interface{})
	Load(bucket, k string) (interface{}, bool)
}
