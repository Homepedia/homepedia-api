package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnection struct {
	Auth *gorm.DB
}

var Connections DBConnection

func Init() {
	var err error
	Connections.Auth, err = initAuthDatabase()
	if err != nil {
		panic("failed to init auth database connection")
	}
}

func initAuthDatabase() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_AUTH")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	var loggerLevel logger.LogLevel
	if os.Getenv("DB_LOG_MODE") == "true" {
		loggerLevel = logger.Info
	} else {
		loggerLevel = logger.Silent
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(loggerLevel),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db, err
}
