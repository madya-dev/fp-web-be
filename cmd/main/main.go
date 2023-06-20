package main

import (
	_ "github.com/joho/godotenv/autoload"
	"hrd-be/pkg/server"
)

func main() {
	//model.InitialMigrate()

	server.StartServer()
}
