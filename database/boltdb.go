package database

import (
	"fmt"

	"github.com/gofrs/uuid"
	bolt "go.etcd.io/bbolt"
)

// BoltDB provides API for using BoltDB
// This function implements Databaser interface.
type BoltDB struct {
	dbpath string
	db     *bolt.DB
}

func (b *BoltDB) open() error {
	var err error
	b.db, err = bolt.Open(b.dbpath, 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	return nil
}

func (b *BoltDB) close() error {
	return b.db.Close()
}

// Initialize initializes BoltDB instance
func (b *BoltDB) Initialize(dbpath string) error {
	b.dbpath = dbpath
	err := b.open()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	return b.close()
}

// Finalize finalizes BoltDB instance
func (b *BoltDB) Finalize() error {
	// nop
	return nil
}

func (b *BoltDB) withDBOpenClose(f func() error) error {
	err := b.open()
	if err != nil {
		return fmt.Errorf("open returned error: %v", err)
	}

	err = f()
	if err != nil {
		return fmt.Errorf("failed to save/load: %v", err)
	}

	return b.close()
}

// Save stores specified key/value to database
func (b *BoltDB) Save(bucket, id string, value []byte) (string, error) {
	var ret string
	return ret, b.withDBOpenClose(
		func() error {
			return b.db.Update(
				func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte(bucket))
					if err != nil {
						return fmt.Errorf("failed to create bucket: %v", err)
					}
					if id == "" {
						u, e := uuid.NewV4()
						if e != nil {
							return fmt.Errorf("failed generate uuid: %v", err)
						}
						id = u.String()
					}

					err = b.Put([]byte(id), value)
					if err != nil {
						return err
					}
					ret = id
					return nil
				})
		})
}

// Load loads data by id
func (b *BoltDB) Load(bucket, id string) ([]byte, error) {
	var ret []byte
	err := b.withDBOpenClose(
		func() error {
			err := b.db.View(
				func(tx *bolt.Tx) error {
					if tx == nil {
						return nil
					}
					b := tx.Bucket([]byte(bucket))
					if b == nil {
						return nil
					}
					v := b.Get([]byte(id))
					if v == nil {
						return nil
					}
					ret = make([]byte, len(v))
					copy(ret, v)
					return nil
				})
			if err != nil {
				return fmt.Errorf("failed load: %v", err)
			}
			return nil
		})
	return ret, err
}

func (b *BoltDB) PhysicalDelete(bucket, id string) error {
	return b.withDBOpenClose(
		func() error {
			err := b.db.Update(
				func(tx *bolt.Tx) error {
					if tx == nil {
						return nil
					}
					b := tx.Bucket([]byte(bucket))
					if b == nil {
						return nil
					}
					if err := b.Delete([]byte(id)); err != nil {
						return fmt.Errorf("failed to physical delete from db: %w", err)
					}
					return nil
				})
			if err != nil {
				return fmt.Errorf("failed to physical delete: %v", err)
			}
			return nil
		})
}

// ForEach loops for all items already saved and invoke specified function
func (b *BoltDB) ForEach(bucket string, f func(k, v []byte) error) error {
	err := b.withDBOpenClose(
		func() error {
			err := b.db.View(
				func(tx *bolt.Tx) error {
					if tx == nil {
						return nil
					}
					b := tx.Bucket([]byte(bucket))
					if b == nil {
						return nil
					}
					err := b.ForEach(f)
					if err != nil {
						return fmt.Errorf("for each stopped: %v", err)
					}
					return nil
				})
			if err != nil {
				return fmt.Errorf("failed load: %v", err)
			}
			return nil
		})
	return err
}
