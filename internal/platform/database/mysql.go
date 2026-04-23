package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MySQL: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Failed to get database instance: %v", err)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
