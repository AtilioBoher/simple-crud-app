package mysqlDB

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AtilioBoher/simple-crud-app/pkg/server/database"
	"github.com/go-sql-driver/mysql"
)

type MySqlDB struct {
	db *sql.DB
}

// NewMySqlDB returns a pointer to an instance of MySqlDB struct, initiated with the mysql.Config
// provided.
func NewMySqlDB(cfg mysql.Config) (*MySqlDB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &MySqlDB{db: db}, nil
}

// CheckConnection checks if the database is up and running.
func (my *MySqlDB) CheckConnection() error {
	return my.db.Ping()
}

func (my *MySqlDB) Get(ctx context.Context, name string) *database.User {
	// An album to hold data from the returned row.
	var (
		user database.User
		id   int
	)

	row := my.db.QueryRow("SELECT * FROM users WHERE user_name = ?;", name)
	if err := row.Scan(&id, &user.Name, &user.Email, &user.Age); err != nil {
		return nil
	}
	return &user
}

func (my *MySqlDB) Create(ctx context.Context, user database.User) error {
	_, err := my.db.Exec("INSERT INTO users (user_name, user_email, user_age) VALUES (?, ?, ?);", user.Name, user.Email, user.Age)
	if err != nil {
		return fmt.Errorf("Create user error: %v", err)
	}
	return nil
}


func (my *MySqlDB) Update(ctx context.Context, user database.User) (*database.User, error) {
	// Check for missing fields
	current := my.Get(ctx, user.Name)
	if user.Email == "" {
		user.Email = current.Email
	}
	if user.Age == 0 {
		user.Age = current.Age
	}

	_, err := my.db.Exec("UPDATE users SET user_email = ?, user_age = ? WHERE user_name = ?; ", user.Email, user.Age, user.Name)
	if err != nil {
		return nil, fmt.Errorf("Update user error: %v", err)
	}
	return &user, nil
}

func (my *MySqlDB) Delete(ctx context.Context, name string) error {
	_, err := my.db.Exec("DELETE FROM users WHERE user_name = ?; ", name)
	if err != nil {
		return fmt.Errorf("Delete user error: %v", err)
	}
	return nil
}