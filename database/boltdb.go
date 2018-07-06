package database

import (
	"errors"
	"strconv"

	bolt "github.com/coreos/bbolt"
)

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

func (b *BoltDB) Initialize(dbpath string) error {
	b.dbpath = dbpath
	err := b.open()
	if err != nil {
		return errors.New("failed to open database: " + err.Error())
	}
	return b.db.Close()
}

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

	return b.db.Close()
}

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
