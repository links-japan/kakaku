package store

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect(dsn string) error {
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = conn

	conn.AutoMigrate(&Asset{})
	return nil
}

func Conn() *gorm.DB {
	return db
}
