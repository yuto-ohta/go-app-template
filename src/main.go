package main

import (
	"go-app-template/src/config"
	"go-app-template/src/config/route"
)

func main() {
	// load config
	config.LoadConfig()

	// server run
	router := route.NewRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
