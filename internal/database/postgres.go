package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect(_ string) (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSLMODE")

	log.Println("DB_HOST =", host)
	log.Println("DB_PORT =", port)
	log.Println("DB_NAME =", name)
	log.Println("DB_USER =", user)
	log.Println("DB_SSLMODE =", ssl)

	if host == "" || user == "" || name == "" {
		return nil, fmt.Errorf("missing database environment variables")
	}

	if port == "" {
		port = "5432"
	}
	if ssl == "" {
		ssl = "require"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, name, ssl,
	)

	log.Println("DB DSN =", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, db.Ping()
}
