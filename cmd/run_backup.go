package cmd

import (
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/service"
	"github.com/spf13/cobra"
)

func RunBackup(version application.VersionConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the backup",
		Long:  "This command will launch the backup.",
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
			backupIds, err := service.Backup(app)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Backups created: %v", backupIds)
		},
	}
}
