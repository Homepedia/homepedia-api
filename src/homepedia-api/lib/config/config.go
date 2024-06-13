package config

import "gorm.io/gorm"

type DBConnection struct {
	Auth *gorm.DB
}

var Connections DBConnection

func Init() {
	var err error
	Connections.Auth, err = InitAuthDatabase()
	if err != nil {
		panic("failed to init auth database connection")
	}
	InitCache()
}
