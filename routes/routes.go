package routes

import (
	"github.com/gin-gonic/gin"
	accountController "hrd-be/internal/account/controller"
	cisController "hrd-be/internal/cis/controller"
	employeeController "hrd-be/internal/employee/controller"
	projectController "hrd-be/internal/project/controller"
	slipController "hrd-be/internal/slip/controller"
	"hrd-be/pkg/jwt"
)

func AccountRoutes(g *gin.RouterGroup) {
	g.POST("/login", accountController.LoginHandler())
	g.POST("/create", jwt.ValidationMiddleware(0), accountController.CreateAccountHandler())
	g.PUT("/", jwt.ValidationMiddleware(1), accountController.EditPasswordHandler())
	g.DELETE("/:username", jwt.ValidationMiddleware(0), accountController.DeleteAccountHandler())
}

func CisRoutes(g *gin.RouterGroup) {
	g.POST("/new", jwt.ValidationMiddleware(1), cisController.NewCisHandler())
	g.GET("/", jwt.ValidationMiddleware(1), cisController.GetAllCisHandler())
	g.GET("/:cis_id", jwt.ValidationMiddleware(1), cisController.CisDetailHandler())
	g.PUT("/:cis_id", jwt.ValidationMiddleware(0), cisController.EditCisHandler())
	g.DELETE("/:cis_id", jwt.ValidationMiddleware(1), cisController.DeleteCisHandler())
}

func EmployeeRoutes(g *gin.RouterGroup) {
	g.GET("/", jwt.ValidationMiddleware(0), employeeController.GetAllEmployeeHandler())
	g.GET("/:employee_id", jwt.ValidationMiddleware(1), employeeController.GetEmployeeDetail())
	g.PUT("/:employee_id", jwt.ValidationMiddleware(0), employeeController.EditEmployeeHandler())
}

func ProjectRoutes(g *gin.RouterGroup) {
	g.POST("/new", jwt.ValidationMiddleware(0), projectController.NewProjectHandler())
	g.GET("/", jwt.ValidationMiddleware(1), projectController.GetAllProjectHandler())
	g.GET("/:project_id", jwt.ValidationMiddleware(1), projectController.GetProjectDetailHandler())
	g.PUT("/:project_id", jwt.ValidationMiddleware(0), projectController.EditProjectHandler())
	g.DELETE("/:project_id", jwt.ValidationMiddleware(0), projectController.DeleteProjectHandler())
}

func SlipRoutes(g *gin.RouterGroup) {
	g.POST("/generate", jwt.ValidationMiddleware(0), slipController.GenerateSlipHandler())
}
