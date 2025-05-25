package http

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/service"
)

func ListBackup(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		backups, err := service.BackupList(app)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, backups)
	}
}

func CreateBackup(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			backupIds, err := service.Backup(app)
			if err != nil {
				log.Printf("%s", err)
				return
			}
			log.Printf("Backup created with ID: %v", backupIds)
		}()
		c.JSON(200, gin.H{"Message": "Backup started"})
	}
}

func GenerateDownloadUrl(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		backupId := c.Param("id")
		url, err := service.GenerateDownloadUrl(app, backupId)
		if err != nil {
			c.JSON(500, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Url": url})
	}
}

func DownloadFile(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		driveFileId := c.Param("id")
		output, err := service.Download(app, driveFileId)
		if err != nil {
			c.JSON(500, gin.H{"Error": err.Error()})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(output.Filename))
		c.Header("Content-Type", output.MimeType)
		c.Data(200, output.MimeType, output.Byte)
	}
}
