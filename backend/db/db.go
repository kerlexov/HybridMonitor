package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host     string
	Username string
	Password string
	Table    string
	Port     int
}

func Init(config DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Host, config.Username, config.Password, config.Table, config.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
