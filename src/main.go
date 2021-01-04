package main

import (
	"go-app-template/src/config/route"
)

func main() {
	// server run
	router := route.NewRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
