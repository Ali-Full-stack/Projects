package middleware

import (
	"net/http"
	redisdb "restaurant-service/inernal/redisDB"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func EmployeePasswordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Next()
			return
		}
		redisClient := redisdb.ConnectRedis()
		id := c.GetHeader("id")
		password := c.GetHeader("password")

		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing Id !",
			})
			return
		}
		empInfo, err := redisClient.VerifyEmployeeAndReturnInfo(id, password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid  ID or Password !",
			})
			return
		}
		enforcer, err := casbin.NewEnforcer("./auth/casbin/auth.conf", "./auth/casbin/auth.csv")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong !",
			})
			return
		}

		ok, _ := enforcer.Enforce("staff", c.Request.URL.Path, c.Request.Method)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Unauthorized access",
			})
			return
		}
		c.Set("role", empInfo.Role)
		c.Set("name", empInfo.Name)
		c.Next()
	}
}
func AdminPasswordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
			c.Next()
			return
		}
		redisClient := redisdb.ConnectRedis()
		role := c.GetHeader("role")
		id := c.GetHeader("id")
		password := c.GetHeader("password")

		if role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing Role !",
			})
			return
		}
		role = strings.ToLower(role)
		if role == "admin" {
			adminInfo, err := redisClient.VerifyEmployeeAndReturnInfo(id, password)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Admin Password is incorrect !",
				})
				return
			}
			c.Set("role", adminInfo.Role)
			c.Set("name", adminInfo.Name)
		}

		enforcer, err := casbin.NewEnforcer("./casbin/auth.conf", "./casbin/auth.csv")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong !",
			})
			return
		}

		ok, _ := enforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Unauthorized access",
			})
			return
		}
		c.Next()
	}
}
