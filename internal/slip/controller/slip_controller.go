package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/internal/slip/dto"
	"hrd-be/internal/slip/service"
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
		defer database.Close(db)

		var account model.Account
		db.Preload("Employee.EmployeeStatus").Where("username = ?", generateInput.Username).
			Find(&account)

		var cuti, izin, sakit sql.NullInt64
		stmt := "SELECT SUM(cis_details.end_date - cis_details.start_date) FROM cis_details INNER JOIN cis ON cis_details.id = cis.cis_detail_id WHERE cis.employee_id = ? AND cis.cis_type_id = ? AND cis_details.start_date >= ? AND cis_details.end_date <= ? AND cis.cis_status_id = 3"
		db.Raw(stmt,
			account.EmployeeID, 1, generateInput.StartPeriode, generateInput.EndPeriode).
			Scan(&cuti)
		db.Raw(stmt,
			account.EmployeeID, 2, generateInput.StartPeriode, generateInput.EndPeriode).
			Scan(&izin)
		db.Raw(stmt,
			account.EmployeeID, 3, generateInput.StartPeriode, generateInput.EndPeriode).
			Scan(&sakit)

		var paidLeave, permission, insurance float64
		var salaryCuts []model.SalaryCut
		db.Find(&salaryCuts)

		paidLeave = 0
		permission = 0
		insurance = 0
		for _, salaryCut := range salaryCuts {
			if salaryCut.Name == "paid_leave" {
				paidLeave = salaryCut.SalaryCut
			}
			if salaryCut.Name == "permission" {
				permission = salaryCut.SalaryCut
			}
			if salaryCut.Name == "insurance" {
				insurance = salaryCut.SalaryCut
			}
		}

		startPeriode, _ := time.Parse("2006-01-02", generateInput.StartPeriode)
		endPeriode, _ := time.Parse("2006-01-02", generateInput.EndPeriode)

		var slip model.SalarySlip
		slip.Name = account.Employee.Name
		slip.Position = account.Employee.Position
		slip.Status = account.Employee.EmployeeStatus.Name
		slip.StartPeriode = startPeriode
		slip.EndPeriode = endPeriode
		slip.BasicSalary = account.Employee.Salary
		slip.Bonus = generateInput.Bonus
		slip.PaidLeave = paidLeave
		slip.Permission = permission
		slip.Insurance = insurance

		var filename string
		err := db.Transaction(func(tx *gorm.DB) error {
			err := tx.Create(&slip).Error
			if err != nil {
				return err
			}

			filename, err = service.GenerateSlip(slip)
			if err != nil {
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

		protocol := "http://"
		if c.Request.TLS != nil {
			protocol = "https://"
		}
		file := protocol + c.Request.Host + "/slips/" + filename
		response.DefaultOK()
		response.Message = "slip generated successfully"
		response.Data = map[string]string{
			"slip": file,
		}
		c.JSON(response.Code, response)
	}
}
