package cmd

import (
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/service"
	"github.com/spf13/cobra"
)

func RetryBackup(version application.VersionConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "retry [id]",
		Short: "Retry a failed backup",
		Long:  "This command will retry a failed backup.",
		Run: func(cmd *cobra.Command, args []string) {
			configFile, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal(err)
			}

			if len(args) < 1 {
				log.Fatal("Backup ID is required for retry")
			}
			app, err := CreateAppFromYml(configFile)
			if err != nil {
				log.Fatalf("Error creating app from config => %v", err)
			}
			app.Mode = application.APP_MODE_CLI
			app.Version = version
			backupId := args[0]
			err = service.BackupRetry(app, backupId)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Backup retry completed for ID: %s", backupId)
		},
	}
}
