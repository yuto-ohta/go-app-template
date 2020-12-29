package main

import (
	"go-app-template/config"
	"go-app-template/config/routes"
)

func main() {
	// load config
	config.LoadConfig()

	// server run
	router := routes.NewRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
