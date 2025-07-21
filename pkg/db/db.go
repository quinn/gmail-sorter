package db

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/api/gmail/v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB holder for DB junk
type db struct {
	db   *bolt.DB
	gorm *gorm.DB
}

var DB *db

func init() {
	path := "./bolt.db"
	d, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatalf("failed to open boltdb: %v", err)
	}

	gormdb, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite db: %v", err)
	}

	if err := gormdb.AutoMigrate(&OAuthAccount{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	DB = &db{db: d, gorm: gormdb}
}

// Upsert insert or update based on key and bucket
func (d *db) Upsert(bucket string, key string, value []byte) error {
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

func (d *db) bucket(tx *bolt.Tx, bucket string) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte(bucket))
	if b != nil {
		return b, nil
	}

	b, err := tx.CreateBucket([]byte(bucket))
	if err != nil {
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return b, nil
}

// Get a key from a bucket
func (d *db) Get(bucket string, key string) (bytes []byte, err error) {
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
func (d *db) Close() error {
	return d.db.Close()
}

// GetAll retrieve all of the objects for the given bucket name
func (d *db) GetAll(bucket string) (objects [][]byte, err error) {
	err = d.db.Update(func(tx *bolt.Tx) (err error) {
		b, err := d.bucket(tx, bucket)

		if err != nil {
			return err
		}

		err = b.ForEach(func(k []byte, v []byte) (err error) {
			objects = append(objects, v)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (db *db) filterKey(accountID uint) string {
	return "filters-" + strconv.FormatUint(uint64(accountID), 10)
}

func (db *db) labelKey(accountID uint) string {
	return "labels-" + strconv.FormatUint(uint64(accountID), 10)
}

func (db *db) Filters(accountID uint) ([]*gmail.Filter, error) {
	filters, err := db.GetAll(db.filterKey(accountID))

	if err != nil {
		return nil, fmt.Errorf("failed to get all filters: %w", err)
	}

	var result []*gmail.Filter

	for _, bytes := range filters {
		var filter gmail.Filter
		if err := json.Unmarshal(bytes, &filter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filter: %w", err)
		}

		result = append(result, &filter)
	}

	return result, nil
}

func (db *db) Label(accountID uint, id string) (*gmail.Label, error) {
	bytes, err := db.Get(db.labelKey(accountID), id)

	if err != nil {
		return nil, err
	}

	if bytes == nil {
		return nil, fmt.Errorf("label %s not found", id)
	}

	var label gmail.Label
	if err = json.Unmarshal(bytes, &label); err != nil {
		return nil, err
	}

	return &label, nil
}

func (db *db) UpsertFilters(accountID uint, filters []*gmail.Filter) error {
	for _, filter := range filters {
		d, err := json.Marshal(filter)
		if err != nil {
			return err
		}

		if err := db.Upsert(db.filterKey(accountID), filter.Id, d); err != nil {
			return err
		}
	}

	return nil
}

func (db *db) UpsertLabels(accountID uint, labels []*gmail.Label) error {
	for _, label := range labels {
		d, err := json.Marshal(label)
		if err != nil {
			return err
		}

		if err := db.Upsert(db.labelKey(accountID), label.Id, d); err != nil {
			return err
		}
	}

	return nil
}

func (db *db) AllFilters(accountID uint) ([]*gmail.Filter, error) {
	filters, err := db.GetAll(db.filterKey(accountID))
	if err != nil {
		return nil, fmt.Errorf("failed to get all filters: %w", err)
	}

	var result []*gmail.Filter

	for _, bytes := range filters {
		var filter gmail.Filter
		if err := json.Unmarshal(bytes, &filter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal filter: %w", err)
		}

		result = append(result, &filter)
	}

	return result, nil
}

func (db *db) UpsertFilter(accountID uint, filter *gmail.Filter) error {
	d, err := json.Marshal(filter)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := db.Upsert(db.filterKey(accountID), filter.Id, d); err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}

	return nil
}
