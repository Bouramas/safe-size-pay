package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	DB *sql.DB
}

func newConfig() (cfg *config, err error) {

	cfg = new(config)
	cfg.DB, err = getMySQL()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func getMySQL() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")

	println("dsn: ", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return db, nil
}
