package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	// Read from env (set in container or shell)
	user := os.Getenv("ADMIN_USER")
	pass := os.Getenv("ADMIN_PASS")

	return func(c *gin.Context) {
		givenUser, givenPass, ok := c.Request.BasicAuth()
		if !ok || givenUser != user || givenPass != pass {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
