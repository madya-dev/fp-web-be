package server

import (
	"github.com/gin-gonic/gin"
	"hrd-be/routes"
	"log"
)

func StartServer() {
	log.Print("INFO StartServer: server is starting")
	router := gin.Default()

	account := router.Group("/account")
	routes.AccountRoutes(account)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("ERROR StartServer fatal error: %v", err)
	}
	log.Println("INFO StartServer: server started successfully")
}
