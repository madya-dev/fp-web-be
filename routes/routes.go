package routes

import (
	"github.com/gin-gonic/gin"
	"hrd-be/internal/account/controller"
)

func AccountRoutes(g *gin.RouterGroup) {
	g.POST("/login", controller.LoginHandler())
	g.POST("/create", controller.CreateAccountHandler())
	g.PUT("/:username", controller.EditPasswordHandler())
	g.DELETE("/:username", controller.DeleteAccountHandler())
}
