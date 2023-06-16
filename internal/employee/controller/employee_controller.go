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
