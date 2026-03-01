package storage

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func OpenMySQL() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")      // environment var
	cfg.Passwd = os.Getenv("DBPASS")    // environment var
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DBADDR")      // e.g. 127.0.0.1:3306
	cfg.DBName = os.Getenv("DBNAME")    // e.g. proxydb
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
		db.Close()
		return nil, err
	}
	return db, nil
}
