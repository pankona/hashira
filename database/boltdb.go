package database

import (
	"errors"
	"strconv"

	bolt "github.com/coreos/bbolt"
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
		return errors.New("failed to open database: " + err.Error())
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
		return errors.New("failed to open database: " + err.Error())
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
		return errors.New("open returned error: " + err.Error())
	}

	err = f()
	if err != nil {
		return errors.New("failed to save/load: " + err.Error())
	}

	return b.close()
}

// Save stores specified key/value to database
func (b *BoltDB) Save(id string, value []byte) error {
	return b.withDBOpenClose(
		func() error {
			return b.db.Update(
				func(tx *bolt.Tx) error {
					b, err := tx.CreateBucketIfNotExists([]byte("taskBucket"))
					if err != nil {
						return errors.New("failed to create bucket: " + err.Error())
					}
					if id == "" {
						n, err := b.NextSequence()
						if err != nil {
							return errors.New("failed to get next sequence from bucket: " + err.Error())
						}
						id = strconv.FormatUint(n, 10)
					}

					return b.Put([]byte(id), value)
				})

		})
}

// Load loads data by id
func (b *BoltDB) Load(id string) ([]byte, error) {
	var ret []byte
	err := b.withDBOpenClose(
		func() error {
			err := b.db.View(
				func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("taskBucket"))
					v := b.Get([]byte(id))
					ret = make([]byte, len(v))
					copy(ret, v)
					return nil
				})
			if err != nil {
				return errors.New("failed load: " + err.Error())
			}
			return nil
		})
	return ret, err
}

// ForEach loops for all items already saved and invoke specified function
func (b *BoltDB) ForEach(f func(k, v []byte) error) error {
	err := b.withDBOpenClose(
		func() error {
			err := b.db.View(
				func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("taskBucket"))
					err := b.ForEach(f)
					if err != nil {
						return errors.New("for each stopped: " + err.Error())
					}
					return nil
				})
			if err != nil {
				return errors.New("failed load: " + err.Error())
			}
			return nil
		})
	return err
}
