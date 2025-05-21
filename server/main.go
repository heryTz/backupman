package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/herytz/backupman/config"
	"github.com/herytz/backupman/core"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.yml", "Path to the config file")
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	config, err := config.YmlToAppConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	app := core.NewApp(config)
	app.Mode = core.APP_MODE_WEB

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
	err = router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}
