package main

import "go-app-template/config"

func main() {
	// server run
	router := config.NewRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
