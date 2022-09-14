package middleware

import (
	"net/http"

	util "github.com/Christomesh/pugasell/utils"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("token")
		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
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

func AuthorizedPermission(c *gin.Context, role string) (err error) {
	resourceRole := c.GetString("role")
	err = nil
	if resourceRole != role {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this route"})
		return
	}
	return err
}

func CheckPermission(c *gin.Context, requestUserId string) (err error) {
	role := c.GetString("role")
	resourceUserId := c.GetString("userId")
	err = nil
	if role == "USER" && resourceUserId != requestUserId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized to access this route"})
		return
	}
	err = AuthorizedPermission(c, role)
	return err
}
