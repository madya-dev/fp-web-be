package server

import (
	"github.com/gin-gonic/gin"
	"hrd-be/routes"
	"log"
)

func StartServer() {
	log.Print("INFO StartServer: server is starting")
	router := gin.Default()

	public := router.Group("/")
	routes.PublicRoutes(public)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("ERROR StartServer fatal error: %v", err)
	}
	log.Println("INFO StartServer: server started successfully")
}
