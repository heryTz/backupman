package http

import (
	"github.com/gin-gonic/gin"
	"github.com/herytz/backupman/core/application"
)

func Auth(app *application.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-Api-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		for _, key := range app.Http.ApiKeys {
			if apiKey == key {
				c.Next()
				return
			}
		}
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
