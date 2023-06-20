package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hrd-be/internal/employee/dto"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/model"
	"hrd-be/pkg/database"
	inputValidator "hrd-be/pkg/validator"
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
		var account model.Account

		db := database.Connection()
		result := db.Preload("Employee.EmployeeStatus").
			Where("employee_id = ?", employeeId).
			Find(&account)

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
			Status     int     `json:"status"`
			Username   string  `json:"username"`
			Role       int     `json:"role"`
		}
		cleanEmployee := CleanEmployee{
			EmployeeID: account.EmployeeID,
			Name:       account.Employee.Name,
			Age:        account.Employee.Age,
			Salary:     account.Employee.Salary,
			Position:   account.Employee.Position,
			Status:     account.Employee.EmployeeStatusID,
			Username:   account.Username,
			Role:       account.Role,
		}

		response.DefaultOK()
		response.Message = "success get employee detail"
		response.Data = map[string]CleanEmployee{
			"employee": cleanEmployee,
		}
		c.JSON(response.Code, response)
	}
}

func EditEmployeeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		employeeId := c.Param("employee_id")
		var editInput dto.EditInput
		if err := c.BindJSON(&editInput); err != nil {
			response.DefaultInternalError()
			response.Data = map[string]string{"errors": err.Error()}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(editInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()

		var count int64
		if db.Where("employee_id = ?", employeeId).First(&model.Account{}).Count(&count); count != 1 {
			response.DefaultNotFound()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			var account model.Account
			tx.Where("employee_id = ?", employeeId).First(&account)
			account.Username = editInput.Username
			account.Role = editInput.Role
			if err := tx.Save(&account).Error; err != nil {
				return err
			}

			var employee model.Employee
			tx.Where("id = ?", employeeId).First(&employee)
			employee.Name = editInput.Name
			employee.Age = editInput.Age
			employee.Salary = editInput.Salary
			employee.Position = editInput.Position
			employee.EmployeeStatusID = editInput.Status
			if err := tx.Save(&employee).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			response.DefaultInternalError()
			response.Data = map[string]string{"errors": err.Error()}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultOK()
		response.Message = "employee detail updated successfully"
		c.JSON(response.Code, response)
	}
}
