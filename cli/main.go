package main

import (
	"log"

	"github.com/herytz/backupman/config"
	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/service"
	"github.com/spf13/cobra"
)

func main() {
	var configFile string
	var app *core.App

	rootCmd := &cobra.Command{
		Use:   "backupman",
		Short: "Backup Manager CLI",
		Long:  "A command line tool for managing backups.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config, err := config.YmlToAppConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}
			app = core.NewApp(config)
			app.Mode = core.APP_MODE_CLI
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./config.yml", "Path to the config file")

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the backup",
		Long:  "This command will launch the backup.",
		Run: func(cmd *cobra.Command, args []string) {
			backupIds, err := service.Backup(app)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Backups created: %v", backupIds)

		},
	}

	retryCmd := &cobra.Command{
		Use:   "retry [id]",
		Short: "Retry a failed backup",
		Long:  "This command will retry a failed backup.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Backup ID is required for retry")
			}
			backupId := args[0]
			err := service.BackupRetry(app, backupId)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Backup retry completed for ID: %s", backupId)
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(retryCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
