package main

import (
	"go-app-template/src/config/route"
)

func main() {
	// run server
	router := route.NewRouter()
	router.Logger.Fatal(router.Start(":1323"))
}
