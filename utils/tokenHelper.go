package utils

import (
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email   string
	User_id string
	Role    string
	jwt.StandardClaims
}

func GenerateToken(c *gin.Context, email string, user_id string, role string) {
	claims := &SignedDetails{
		Email:   email,
		User_id: user_id,
		Role:    role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err.Error())
		return
	}
	c.SetCookie("token", token, 60, "/", "localhost", false, true)

}
