package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/internal/slip/dto"
	"hrd-be/model"
	"hrd-be/pkg/database"
	inputValidator "hrd-be/pkg/validator"
	"time"
)

func GenerateSlipHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		var generateInput dto.GenerateInput
		if err := c.Bind(&generateInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(generateInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()

		var employee model.Employee
		db.Preload("EmployeeStatus").Where("id = ?", generateInput.EmployeeID).
			Find(&employee)

		var cuti, izin, sakit sql.NullInt64
		db.Raw("SELECT SUM(cis_details.end_date - cis_details.start_date) FROM cis_details INNER JOIN cis ON cis_details.id = cis.cis_detail_id WHERE cis.employee_id = ? AND cis.cis_type_id = ? AND cis_details.start_date >= ? AND cis_details.end_date <= ? AND cis.cis_status_id = 2", generateInput.EmployeeID, 1, generateInput.StartPeriode, generateInput.EndPeriode).
			Scan(&cuti)
		db.Raw("SELECT SUM(cis_details.end_date - cis_details.start_date) FROM cis_details INNER JOIN cis ON cis_details.id = cis.cis_detail_id WHERE cis.employee_id = ? AND cis.cis_type_id = ? AND cis_details.start_date >= ? AND cis_details.end_date <= ? AND cis.cis_status_id = 2", generateInput.EmployeeID, 2, generateInput.StartPeriode, generateInput.EndPeriode).
			Scan(&izin)
		db.Raw("SELECT SUM(cis_details.end_date - cis_details.start_date) FROM cis_details INNER JOIN cis ON cis_details.id = cis.cis_detail_id WHERE cis.employee_id = ? AND cis.cis_type_id = ? AND cis_details.start_date >= ? AND cis_details.end_date <= ? AND cis.cis_status_id = 2", generateInput.EmployeeID, 3, generateInput.StartPeriode, generateInput.EndPeriode).
			Scan(&sakit)

		var paidLeave, permission, insurance float64
		var salaryCuts model.SalaryCut
		db.Where("name = ?", "paid_leave").Find(&salaryCuts)
		paidLeave = float64(cuti.Int64) * salaryCuts.SalaryCut

		db.Where("name = ?", "permission").Find(&salaryCuts)
		permission = float64(izin.Int64) * salaryCuts.SalaryCut

		db.Where("name = ?", "insurance").Find(&salaryCuts)
		insurance = salaryCuts.SalaryCut

		startPeriode, _ := time.Parse("2006-01-02", generateInput.StartPeriode)
		endPeriode, _ := time.Parse("2006-01-02", generateInput.EndPeriode)

		var slip model.SalarySlip
		slip.Name = employee.Name
		slip.Position = employee.Position
		slip.Status = employee.EmployeeStatus.Name
		slip.StartPeriode = startPeriode
		slip.EndPeriode = endPeriode
		slip.BasicSalary = employee.Salary
		slip.Bonus = generateInput.Bonus
		slip.PaidLeave = paidLeave
		slip.Permission = permission
		slip.Insurance = insurance

		result := db.Create(&slip)
		if result.Error != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultOK()
		response.Message = "slip generated successfully"
		response.Data = map[string]model.SalarySlip{
			"slip": slip,
		}
		c.JSON(response.Code, response)
	}
}
