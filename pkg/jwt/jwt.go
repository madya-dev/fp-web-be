package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/model"
	"hrd-be/pkg/database"
	"log"
	"os"
	"strings"
	"time"
)

type CustomClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

func CreateToken(id int, username string, role int) string {
	mySigningKey := []byte(os.Getenv("SIGNKEY"))
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &CustomClaims{
		id,
		username,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatalf("ERROR CreateToken fatal error: %v", err)
	}

	return tokenString
}

func ValidationMiddleware(role int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		mySigningKey := []byte(os.Getenv("SIGNKEY"))
		authorization := c.GetHeader("Authorization")
		response.DefaultUnauthorized()
		if authorization == "" {
			response.Data = map[string]string{"error": "missing Authorization request header"}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		authorizationMap := strings.Split(authorization, " ")
		tokenString := authorizationMap[len(authorizationMap)-1]

		claims := &CustomClaims{}

		// parse token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				response.Data = map[string]string{"error": "invalid token"}
				c.AbortWithStatusJSON(response.Code, response)
				return
			}

			response.DefaultForbidden()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		// token validation
		response.DefaultForbidden()
		response.Data = map[string]string{"error": "invalid token"}
		if !token.Valid {
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		// metadata validation
		db := database.Connection()

		var account model.Account
		err = db.Where("username = ?", claims.Username).Find(&account).Error
		if err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		if account.Username != claims.Username || account.Role != claims.Role || claims.Role > role {
			response.DefaultForbidden()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		// creating refresh token
		if time.Until(claims.ExpiresAt.Time) < 10*time.Minute {
			expirationTime := time.Now().Add(1 * time.Hour)
			claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
		}

		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString(mySigningKey)
		if err != nil {
			response.DefaultInternalError()
			response.Data = map[string]string{"error": "refresh token creation error"}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		c.Header("Authorization", tokenString)
		c.Set("claims", claims)
		c.Next()
	}
}
