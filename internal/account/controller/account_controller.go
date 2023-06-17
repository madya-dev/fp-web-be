package controller

import (
	"hrd-be/internal/account/dto"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/model"
	"hrd-be/pkg/database"
	"hrd-be/pkg/jwt"
	inputValidator "hrd-be/pkg/validator"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		var loginInput dto.LoginInput
		if err := c.BindJSON(&loginInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(loginInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()
		var account model.Account
		db.Select("username, password, role, employee_id").Where("username = ?", loginInput.Username).Find(&account)

		err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginInput.Password))
		if account.Username != loginInput.Username || err != nil {
			response.DefaultNotFound()
			response.Message = "invalid username or password"
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(os.Getenv("DEFAULT_PASS"))); err == nil {
			response.DefaultUnauthorized()
			response.Data = map[string]string{"error": "user still use default password"}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		jwtString := jwt.CreateToken(account.Username, account.Role)
		response.DefaultOK()
		response.Message = "login success"
		response.Data = map[string]interface{}{
			"username":    account.Username,
			"employee_id": account.EmployeeID,
		}
		c.Header("Authorization", jwtString)
		c.JSON(response.Code, response)
	}
}

func CreateAccountHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		var createInput dto.CreateInput

		if err := c.BindJSON(&createInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(createInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()
		err := db.Transaction(func(tx *gorm.DB) error {
			employee := model.Employee{}
			if err := tx.Create(&employee).Error; err != nil {
				return err
			}

			bcryptPass, _ := bcrypt.GenerateFromPassword([]byte(os.Getenv("DEFAULT_PASS")), 10)
			accounts := model.Account{
				Username:   createInput.Username,
				Email:      createInput.Email,
				Password:   string(bcryptPass),
				Role:       createInput.Role,
				EmployeeID: employee.ID,
			}
			if err := tx.Create(&accounts).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			response.DefaultConflict()
			response.Data = map[string]string{
				"errors": err.Error(),
			}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultCreated()
		response.Message = "account created successfully"
		c.JSON(response.Code, response)
	}
}

func EditPasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var response globalResponse.Response
		var editPasswordInput dto.EditPasswordInput
		if err := c.BindJSON(&editPasswordInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(editPasswordInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()
		bcryptPass, _ := bcrypt.GenerateFromPassword([]byte(editPasswordInput.Password), 10)
		result := db.Model(model.Account{}).
			Where("username = ? AND email = ?", username, editPasswordInput.Email).
			Update("password", string(bcryptPass))
		if result.RowsAffected == 0 {
			response.DefaultNotFound()
			response.Data = map[string]string{
				"errors": "invalid username or email verification",
			}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultOK()
		response.Message = "password updated successfully"
		c.JSON(response.Code, response)
	}
}

func DeleteAccountHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var response globalResponse.Response

		db := database.Connection()
		var account model.Account
		var count int64
		result := db.Select("employee_id").Where("username = ?", username).Find(&account)
		if result.Count(&count); count == 0 {
			response.DefaultNotFound()
			response.Data = map[string]string{
				"errors": "username not found",
			}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("username = ? AND employee_id = ?", username, account.EmployeeID).
				Delete(&model.Account{}).
				Error; err != nil {
				return err
			}

			if err := tx.Where("id = ?", account.EmployeeID).
				Delete(&model.Employee{}).
				Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultOK()
		response.Message = "account deleted successfully"
		c.JSON(response.Code, response)
	}
}
