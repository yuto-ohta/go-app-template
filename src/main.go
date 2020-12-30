package main

import (
	"go-app-template/src/config"
	"go-app-template/src/config/routes"
)

func main() {
	// load config
	config.LoadConfig()

	// server run
	router := routes.NewRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
