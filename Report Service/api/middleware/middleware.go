package middleware

import (
	"log"
	"net/http"
	"rabbitmq-topic/auth"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func VerifyTokenMIddleware(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/swagger/") {
		c.Next()
		return
	}

	token := c.GetHeader("token")
	
	if token == "" {
		c.AbortWithStatusJSON(403, gin.H{
			"error": "Missing Token",
		})
		return
	}
	
	role, err :=auth.VerifyToken(token)
	if err != nil {
		log.Println("verify token :",err)
		c.AbortWithStatusJSON(403, gin.H{
			"status" : "Permission Denied !",
		})
	}

	enforcer, err := casbin.NewEnforcer("./auth/casbin/auth.conf", "./auth/casbin/auth.csv")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong !",
		})
		return
	}

	ok, err:= enforcer.Enforce(role, c.Request.URL.String(), c.Request.Method)
	if !ok {
		log.Println("enforce:", err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": "Permission Denied",
		})
		return
	}

	c.Next()
}
