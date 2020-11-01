package kakaku

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func Connect() error {
	conn, err := gorm.Open(mysql.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
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
