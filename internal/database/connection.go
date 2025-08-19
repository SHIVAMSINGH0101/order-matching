package database

import (
	"database/sql"
	"fmt"

	"github.com/SHIVAMSINGH0101/go-demo/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}
