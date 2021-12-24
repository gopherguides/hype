package hype

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/boltdb/bolt"
)

// ErrCacheMiss is returned when a cache miss occurs.
type ErrCacheMiss struct {
	Key string
}

func (ecm ErrCacheMiss) Error() string {
	return fmt.Sprintf("cache miss: %s", ecm.Key)
}

// Cache is a cache for storing and retrieving data.
type Cache struct {
	Root string // default: pwd

	sync.RWMutex
	db *bolt.DB
}

// Open the cache at the given path.
func (c *Cache) Open() error {
	if c == nil {
		return fmt.Errorf("cache is nil")
	}

	os.MkdirAll(c.Root, 0755)

	fp := filepath.Join(c.Root, "cache.db")
	db, err := bolt.Open(fp, 0755, nil)
	if err != nil {
		return err
	}

	c.Lock()
	c.db = db
	c.Unlock()

	return nil
}

// Close the cache.
func (c *Cache) Close() error {
	if c == nil {
		return nil
	}

	c.Lock()
	if c.db != nil {
		c.db.Close()
	}

	defer c.Unlock()
	return nil
}

// DB returns the bolt database for the given root.
func (c *Cache) DB(root string) (*bolt.DB, error) {
	if c == nil {
		return nil, fmt.Errorf("cache is nil")
	}

	c.Lock()
	if c.db != nil {
		defer c.Unlock()
		err := c.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(c.BucketName())
			return err
		})
		return c.db, err
	}
	c.Unlock()

	db, err := bolt.Open(filepath.Join(root, "hype.db"), 0755, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(c.BucketName())
		return err
	})

	if err != nil {
		return nil, err
	}

	c.Lock()
	c.db = db
	c.Unlock()

	return c.db, nil
}

// BucketName returns the name of the cache bucket.
func (c *Cache) BucketName() []byte {
	return []byte("hype-cache")
}

// Store the given data in the cache.
func (c *Cache) Store(root string, key string, value []byte) error {
	// fmt.Printf("[CACHE]: store\t%q\n", key)
	db, err := c.DB(root)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		buck, err := tx.CreateBucketIfNotExists(c.BucketName())
		if err != nil {
			return err
		}

		return buck.Put([]byte(key), value)
	})
}

// Retrieve the data from the cache.
func (c *Cache) Retrieve(root string, key string) ([]byte, error) {
	// fmt.Printf("[CACHE]: retrieve\t%q\n", key)
	db, err := c.DB(root)
	if err != nil {
		return nil, err
	}

	var res []byte

	err = db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket(c.BucketName())
		if buck == nil {
			return fmt.Errorf("bucket not found: %s", c.BucketName())
		}

		res = buck.Get([]byte(key))
		if res == nil || len(res) == 0 {
			// fmt.Printf("[CACHE]: MISS\t%q\n", key)
			return ErrCacheMiss{Key: key}
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	// fmt.Printf("[CACHE]: HIT\t%q\n", key)
	return res, nil
}
