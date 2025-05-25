package http

import (
	"github.com/gin-gonic/gin"
	"github.com/herytz/backupman/core"
)

func Auth(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-Api-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		for _, key := range app.ApiKeys {
			if apiKey == key {
				c.Next()
				return
			}
		}
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
	}
}
