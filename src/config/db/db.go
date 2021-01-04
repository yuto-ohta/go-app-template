package db

import (
	"fmt"
	"go-app-template/src/config"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Conn *gorm.DB
	err  error
)

const (
	_appEnvKey  = "APP_ENV"
	_envDefault = "development_without_docker"
)

func init() {
	env := os.Getenv(_appEnvKey)

	var params map[interface{}]interface{}
	switch env {
	case "development":
		params = config.GetConfig()[env].(map[interface{}]interface{})
	default:
		params = config.GetConfig()[_envDefault].(map[interface{}]interface{})
	}

	USER := params["user"].(string)
	PASS := params["password"].(string)
	portFrom := params["port_from"].(string)
	portTo := params["port_to"].(int)
	PROTOCOL := fmt.Sprintf("tcp(%v:%v)", portFrom, portTo)
	DBNAME := params["db_name"].(string)
	PARAMS := params["params"].(string)

	DSN := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?" + PARAMS
	Conn, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
}
