package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

type CustomClaims struct {
	Username string
	Role     int
	jwt.RegisteredClaims
}

func CreateToken(username string, role int) string {
	mySigningKey := []byte(os.Getenv("SIGNKEY"))
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &CustomClaims{
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
