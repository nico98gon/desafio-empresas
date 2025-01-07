package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitMySQL() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN") // Formato: user:password@tcp(host:port)/dbname
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error abriendo MySQL: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a MySQL: %w", err)
	}
	return db, nil
}
