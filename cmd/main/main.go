package main

import (
	_ "github.com/joho/godotenv/autoload"
	"hrd-be/pkg/database"
)

func main() {
	database.Connection()
}
