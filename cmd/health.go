package cmd

import (
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/service"
	"github.com/spf13/cobra"
)

func Health(version application.VersionConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Health check",
		Long:  "Perform a health check on the backup system.",
		Run: func(cmd *cobra.Command, args []string) {
			configFile, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal(err)
			}
			app, err := CreateAppFromYml(configFile)
			if err != nil {
				log.Fatalf("Error creating app from config => %v", err)
			}
			app.Mode = application.APP_MODE_CLI
			app.Version = version
			result, err := service.Health(app)
			if err != nil {
				log.Fatal(err)
			}
			message := "Backup system is healthy"
			if result.Status == lib.HEALTH_DOWN {
				message = "Backup system is unhealthy"
			}
			log.Printf("%s. Details: %s", message, result)
		},
	}
}
