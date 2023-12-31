package bolt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/AtilioBoher/simple-crud-app/pkg/server/database"
)

// Bolt is the Bolt database.
// It satisfies the Database interface
type Bolt struct {
	db *bolt.DB
}

const (
	dbName     = "test.db"
	bucketName = "users"
)

// New returns a new Bolt implementation.
func New(directory string) (*Bolt, error) {
	err := os.Mkdir("data", 0750)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	db, err := bolt.Open(fmt.Sprintf("%s/%s", directory, dbName), 0600, nil)
	if err != nil {
		return nil, err
	}

	// Ensure that the bucket exists.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Bolt{
		db: db,
	}, nil
}

type userinfo struct {
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Close closes the database.
// Make sure to close the database once used.
func (b *Bolt) Close(ctx context.Context) {
	b.db.Close()
}

// Create implements the Database interface.
func (b *Bolt) Create(ctx context.Context, user database.User) error {
	userinfo := userinfo{
		Email: user.Email,
		Age:   user.Age,
	}

	v, err := json.Marshal(userinfo)
	if err != nil {
		return err
	}

	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(user.Name), v)
		return err
	})
	return nil
}

// Get implements the Database interface.
func (b *Bolt) Get(ctx context.Context, name string) (user *database.User) {
	var raw []byte
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		raw = b.Get([]byte(name))
		return nil
	})
	if len(raw) == 0 {
		return
	}

	var u database.User

	err := json.Unmarshal(raw, &u)
	if err != nil {
		log.Fatalf("Database corruption: %v", err)
	}
	u.Name = name
	user = &u
	return
}

// Update implements the Database interface.
func (b *Bolt) Update(ctx context.Context, user database.User) (*database.User, error) {
	var raw []byte
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		raw = b.Get([]byte(user.Name))
		return nil
	})
	var current database.User
	err := json.Unmarshal(raw, &current)
	if err != nil {
		log.Fatalf("Database corruption: %v", err)
	}
	current.Age = user.Age
	current.Email = user.Email
	current.Name = user.Name

	// Write back.
	v, err := json.Marshal(current)
	if err != nil {
		panic(fmt.Sprintf("culd not marshal user: %v", err))
	}
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(user.Name), v)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &current, nil
}

// Delete implements the Database interface.
func (b *Bolt) Delete(ctx context.Context, name string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Delete([]byte(name))
		return err
	})
}
