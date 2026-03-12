package storage

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func mustEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing environment variable %s", key)
	}
	return value, nil
}

func OpenMySQL() (*sql.DB, error) {
	user, err := mustEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := mustEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := mustEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := mustEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbName, err := mustEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = net.JoinHostPort(host, port)
	cfg.DBName = dbName
	cfg.ParseTime = true
	cfg.Params = map[string]string{
		"charset": "utf8mb4",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}