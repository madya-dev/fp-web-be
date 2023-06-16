package controller

import (
	"github.com/gin-gonic/gin"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/internal/project/dto"
	"hrd-be/model"
	"hrd-be/pkg/database"
	inputValidator "hrd-be/pkg/validator"
	"math"
	"strconv"
	"time"
)

func NewProjectHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		var newProjectInput dto.NewProjectInput
		if err := c.BindJSON(&newProjectInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(newProjectInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		startDate, _ := time.Parse("2006-01-02", newProjectInput.StartDate)
		endDate, _ := time.Parse("2006-01-02", newProjectInput.EndDate)

		var employees []model.Employee
		for _, each := range newProjectInput.EmployeesID {
			var employee model.Employee
			employee.ID = each
			employees = append(employees, employee)
		}

		db := database.Connection()
		project := model.Project{
			Name:      newProjectInput.Name,
			Client:    newProjectInput.Client,
			Budget:    newProjectInput.Budget,
			StartDate: startDate,
			EndDate:   endDate,
			Employees: employees,
		}

		if err := db.Create(&project).Error; err != nil {
			response.DefaultInternalError()
			response.Data = map[string]string{"errors": err.Error()}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultCreated()
		response.Message = "project created successfully"
		c.JSON(response.Code, response)
	}
}

func GetAllProjectHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		employeeId := c.Query("employee_id")
		currentPage := c.Query("page")
		currentPageInt, _ := strconv.Atoi(currentPage)
		if currentPageInt < 1 {
			currentPageInt = 1
		}
		perPage := 20
		firstData := (currentPageInt * perPage) - perPage

		type Project struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Client    string `json:"client"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		}

		db := database.Connection()
		var totalData int64
		var totalPage int
		db.Model(&model.Project{}).Count(&totalData)
		totalPage = int(math.Ceil(float64(totalData) / float64(perPage)))

		var projects []model.Project
		query := db.Preload("Employees")
		if employeeId != "" {
			query.Joins("JOIN project_employees ON projects.id = project_employees.project_id").
				Where("employee_id = ?", employeeId)
		}
		query.Order("projects.id DESC").Limit(perPage).Offset(firstData).Find(&projects)

		var cleanProjects []Project
		for _, each := range projects {
			var cleanProject Project
			cleanProject.ID = each.ID
			cleanProject.Name = each.Name
			cleanProject.Client = each.Client
			cleanProject.StartDate = each.StartDate.String()
			cleanProject.EndDate = each.EndDate.String()

			cleanProjects = append(cleanProjects, cleanProject)
		}

		response.DefaultOK()
		response.Message = "success get projects"
		response.Data = map[string]interface{}{
			"cis_list": cleanProjects,
			"pagination": map[string]int{
				"current_page": currentPageInt,
				"total_page":   totalPage,
			},
		}
		c.JSON(response.Code, response)
	}
}

func GetProjectDetailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		projectId := c.Param("project_id")

		var project model.Project
		db := database.Connection()
		result := db.Preload("Employees").
			Where("id = ?", projectId).Find(&project)

		var count int64
		if result.Count(&count); count == 0 {
			response.DefaultNotFound()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		type Employee struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Position string `json:"position"`
		}

		type Project struct {
			ID        int        `json:"id"`
			Name      string     `json:"name"`
			Client    string     `json:"client"`
			Budget    float64    `json:"budget"`
			StartDate string     `json:"start_date"`
			EndDate   string     `json:"end_date"`
			Longtime  int        `json:"longtime"`
			Assign    []Employee `json:"assign"`
		}

		var employees []Employee
		for _, each := range project.Employees {
			var employee Employee
			employee.ID = each.ID
			employee.Name = each.Name
			employee.Position = each.Position

			employees = append(employees, employee)
		}

		cleanProject := Project{
			ID:        project.ID,
			Name:      project.Name,
			Client:    project.Client,
			Budget:    project.Budget,
			StartDate: project.StartDate.String(),
			EndDate:   project.EndDate.String(),
			Longtime:  project.Longtime,
			Assign:    employees,
		}

		response.DefaultOK()
		response.Message = "get project detail success"
		response.Data = map[string]interface{}{
			"project": cleanProject,
		}
		c.JSON(response.Code, response)
	}
}
