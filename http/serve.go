package http

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/herytz/backupman/core/application"
)

func Serve(app *application.App, port int) error {
	app.Mode = application.APP_MODE_WEB

	scheduler, err := SetupScheduler(app)
	if err != nil {
		log.Fatal(err)
	}
	scheduler.Start()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiRouter := router.Group("/api", Auth(app))
	apiRouter.GET("/backups", ListBackup(app))
	apiRouter.POST("/backups", CreateBackup(app))
	apiRouter.GET("/backups/:id/generate-download-url", GenerateDownloadUrl(app))
	apiRouter.GET("/backups/:id/download", DownloadFile(app))

	fmt.Printf("Server is running on port %d\n", port)
	return router.Run(fmt.Sprintf(":%d", port))
}
