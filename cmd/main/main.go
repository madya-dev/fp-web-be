package main

import (
	_ "github.com/joho/godotenv/autoload"
	"hrd-be/db"
	"hrd-be/pkg/server"
)

func main() {
	db.InitialMigrate()

	server.StartServer()
}
