package main

import (
	"log"

	"github.com/herytz/backupman/cmd"
	"github.com/herytz/backupman/core/application"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commitSHA = "none"
	buildDate = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "backupman",
		Short: "Backup Manager CLI",
		Long:  "A command line tool for managing backups.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringP("config", "c", "./config.yml", "Path to the config file")

	versionConfig := application.VersionConfig{
		Version:   version,
		CommitSHA: commitSHA,
		BuildDate: buildDate,
	}

	rootCmd.AddCommand(cmd.RunBackup(versionConfig))
	rootCmd.AddCommand(cmd.RetryBackup(versionConfig))
	rootCmd.AddCommand(cmd.ServeBackup(versionConfig))
	rootCmd.AddCommand(cmd.Version(versionConfig))
	rootCmd.AddCommand(cmd.Health(versionConfig))
	rootCmd.AddCommand(cmd.AuthGoogle(versionConfig))

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
