package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hrd-be/routes"
	"log"
	"time"
)

func StartServer() {
	log.Print("INFO StartServer: server is starting")
	router := gin.Default()

	config := cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	router.MaxMultipartMemory = 2 << 20
	router.Static("files", "./files")

	account := router.Group("/account")
	routes.AccountRoutes(account)

	cis := router.Group("/cis")
	routes.CisRoutes(cis)

	employee := router.Group("/employee")
	routes.EmployeeRoutes(employee)

	project := router.Group("/project")
	routes.ProjectRoutes(project)

	slip := router.Group("/slip")
	routes.SlipRoutes(slip)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("ERROR StartServer fatal error: %v", err)
	}
	log.Println("INFO StartServer: server started successfully")
}
