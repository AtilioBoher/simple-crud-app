package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AtilioBoher/simple-crud-app/pkg/server"
	"github.com/AtilioBoher/simple-crud-app/pkg/server/database/mysqlDB"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	address := ":8080"
	r := mux.NewRouter()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't be loaded")
	}

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               os.Getenv("MYSQL_ROOT_PASSWORD"),
		Net:                  "tcp",
		Addr:                 "db:3306",
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
	}
	// Get a database handle.

	db, err := mysqlDB.NewMySqlDB(cfg)
	if err != nil {
		log.Fatalf("failed to start database: %v", err)
	}

	srv := server.New(ctx, db)
	r.HandleFunc("/", srv.HandleIndex)
	r.HandleFunc("/users/create", srv.HandleCreateUsers)
	r.HandleFunc("/users/{name}", srv.HandleUsers)
	s := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Start server: %v", address)
	log.Fatal(s.ListenAndServe())
}
