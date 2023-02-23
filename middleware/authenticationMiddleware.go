package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	util "github.com/meshachdamilare/pugasell/utils"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var token string
		// authHeader := c.Request.Header.Get("token")
		// if authHeader == "" && strings.HasPrefix(authHeader, "Bearer") {
		// 	token = strings.Split(authHeader, " ")[1]
		// }

		token, _ = c.Cookie("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication Invalid"})
			c.Abort()
			return
		}
		claims, err := util.ValidateToken(token)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("userId", claims.User_id)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AuthorizedPermissions(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		resourceRole := c.GetString("role")
		if resourceRole != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to access this route"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CheckPermission(c *gin.Context, requestUserId string) (err error) {
	role := c.GetString("role")
	resourceUserId := c.GetString("userId")
	err = nil
	if role == "USER" && resourceUserId != requestUserId {
		err = errors.New("not Authorized to access this resource")
		return err
	}
	return err
}
