package main

import (
	"MyGram/database"
	"MyGram/router"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(":8000")
}
