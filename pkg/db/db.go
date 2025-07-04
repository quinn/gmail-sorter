package db

import (
	bolt "go.etcd.io/bbolt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB holder for DB junk
type DB struct {
	db   *bolt.DB
	gorm *gorm.DB
}

// Upsert insert or update based on key and bucket
func (d *DB) Upsert(bucket string, key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		b, err := d.bucket(tx, bucket)

		if err != nil {
			return
		}

		err = b.Put([]byte(key), value)

		if err != nil {
			return
		}

		return nil
	})
}

func (d *DB) bucket(tx *bolt.Tx, bucket string) (b *bolt.Bucket, err error) {
	b = tx.Bucket([]byte(bucket))

	if b == nil {
		b, err = tx.CreateBucket([]byte(bucket))
	}

	if err != nil {
		return
	}

	return
}

// Get a key from a bucket
func (d *DB) Get(bucket string, key string) (bytes []byte, err error) {
	err = d.db.Update(func(tx *bolt.Tx) (err error) {
		b, err := d.bucket(tx, bucket)

		if err != nil {
			return
		}

		bytes = b.Get([]byte(key))

		return nil
	})

	if err != nil {
		return
	}

	return
}

// Close closes boltdb
func (d *DB) Close() error {
	return d.db.Close()
}

// GetAll retrieve all of the objects for the given bucket name
func (d *DB) GetAll(bucket string) (objects [][]byte, err error) {
	err = d.db.View(func(tx *bolt.Tx) (err error) {
		b, err := d.bucket(tx, bucket)

		if err != nil {
			return
		}

		err = b.ForEach(func(k []byte, v []byte) (err error) {
			objects = append(objects, v)
			return
		})

		if err != nil {
			return
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

// NewDB inits it
func NewDB() (*DB, error) {
	path := "./bolt.db"
	db, err := bolt.Open(path, 0666, nil)

	if err != nil {
		return nil, err
	}

	gormdb, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := gormdb.AutoMigrate(&OAuthAccount{}); err != nil {
		return nil, err
	}

	return &DB{db: db, gorm: gormdb}, nil
}
