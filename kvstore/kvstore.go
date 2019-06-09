package kvstore

type KVStore interface {
	Store(bucket, k string, v interface{})
	Load(bucket, k string) (interface{}, bool)
}
