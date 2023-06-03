package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"hrd-be/internal/account/dto"
	"hrd-be/internal/global/auth"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/model"
	"hrd-be/pkg/database"
	inputValidator "hrd-be/pkg/validator"
	"os"
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
			response.DefaultNotAcceptable()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()
		var account model.Account
		db.Select("username, password, role").Where("username = ?", loginInput.Username).Find(&account)

		err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginInput.Password))
		if account.Username != loginInput.Username || err != nil {
			response.DefaultNotFound()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(os.Getenv("DEFAULT_PASS"))); err == nil {
			response.DefaultUnauthorized()
			response.Data = map[string]string{"error": "user still use default password"}
			fmt.Println(os.Getenv("DEFAULT_PASS"))
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		jwtString := auth.CreateToken(account.Username, account.Role)
		response.DefaultOK()
		response.Message = "login success"
		c.Header("Authorization", jwtString)
		c.JSON(response.Code, response)
	}
}
