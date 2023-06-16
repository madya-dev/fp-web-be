package controller

import (
	"github.com/gin-gonic/gin"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/model"
	"hrd-be/pkg/database"
	"math"
	"strconv"
)

func GetAllEmployeeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		var account []model.Account
		currentPage, _ := strconv.Atoi(c.Query("page"))
		if currentPage < 1 {
			currentPage = 1
		}
		perPage := 20
		firstData := (currentPage * perPage) - perPage

		db := database.Connection()

		var totalData int64
		db.Model(&model.Account{}).Count(&totalData)
		totalPage := int(math.Ceil(float64(totalData) / float64(perPage)))

		db.Preload("Employee").
			Limit(perPage).Offset(firstData).Find(&account)

		type Employee struct {
			Username   string `json:"username"`
			Email      string `json:"email"`
			EmployeeID int    `json:"employee_id"`
			Name       string `json:"name"`
		}
		var cleanEmployee []Employee
		for _, each := range account {
			var clean Employee
			clean.Username = each.Username
			clean.Email = each.Email
			clean.EmployeeID = each.EmployeeID
			clean.Name = each.Employee.Name

			cleanEmployee = append(cleanEmployee, clean)
		}

		response.DefaultOK()
		response.Message = "get employee list success"
		response.Data = map[string]interface{}{
			"employee_list": cleanEmployee,
			"pagination": map[string]int{
				"current_page": currentPage,
				"total_page":   totalPage,
			},
		}
		c.JSON(response.Code, response)
	}
}

func GetEmployeeDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		employeeId := c.Param("employee_id")
		var employee model.Employee

		db := database.Connection()
		result := db.Preload("EmployeeStatus").
			Where("id = ?", employeeId).
			Find(&employee)

		var count int64
		if result.Count(&count); count == 0 {
			response.DefaultNotFound()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		type CleanEmployee struct {
			EmployeeID int     `json:"employee_id"`
			Name       string  `json:"name"`
			Age        int     `json:"age"`
			Salary     float64 `json:"salary"`
			Position   string  `json:"position"`
			Status     string  `json:"status"`
		}
		cleanEmployee := CleanEmployee{
			EmployeeID: employee.ID,
			Name:       employee.Name,
			Age:        employee.Age,
			Salary:     employee.Salary,
			Position:   employee.Position,
			Status:     employee.EmployeeStatus.Name,
		}

		response.DefaultOK()
		response.Message = "success get employee detail"
		response.Data = map[string]CleanEmployee{
			"employee": cleanEmployee,
		}
		c.JSON(response.Code, response)
	}
}
