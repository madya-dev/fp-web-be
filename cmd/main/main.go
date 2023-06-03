package main

import (
	_ "github.com/joho/godotenv/autoload"
	"hrd-be/db"
)

func main() {
	db.InitialMigrate()
}
