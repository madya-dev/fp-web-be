package routes

import (
	"github.com/gin-gonic/gin"
	accountController "hrd-be/internal/account/controller"
	cisController "hrd-be/internal/cis/controller"
)

func AccountRoutes(g *gin.RouterGroup) {
	g.POST("/login", accountController.LoginHandler())
	g.POST("/create", accountController.CreateAccountHandler())
	g.PUT("/:username", accountController.EditPasswordHandler())
	g.DELETE("/:username", accountController.DeleteAccountHandler())
}

func CisRoutes(g *gin.RouterGroup) {
	g.POST("/new", cisController.NewCisHandler())
	g.GET("/", cisController.GetAllCisHandler())
	g.GET("/:cis_id", cisController.CisDetailHandler())
	g.PUT("/:cis_id", cisController.EditCisHandler())
	g.DELETE("/:cis_id", cisController.DeleteCisHandler())
}
