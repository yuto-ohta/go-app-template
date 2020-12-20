package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Conn *gorm.DB
	err  error
)

func init() {
	USER := "root"
	PASS := "mysql"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "development"
	PARAM := "charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true"

	DSN := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?" + PARAM
	Conn, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
}
