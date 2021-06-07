package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CreateConnection() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	fmt.Printf("%s:%s@(%s:%s)/%s\n", user, password, host, port, dbName)
	db, err := gorm.Open(
		"mysql",
		fmt.Sprintf(
			// "user:pass@tcp(ip:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbName,
		),
	)
	db.SingularTable(true)
	return db, err
}
