package main

import (
	"log"

	"github.com/herytz/backupman/cmd"
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

	rootCmd.AddCommand(cmd.RunBackup())
	rootCmd.AddCommand(cmd.RetryBackup())
	rootCmd.AddCommand(cmd.ServeBackup())
	rootCmd.AddCommand(cmd.Version(cmd.VersionParams{
		Version:   version,
		CommitSHA: commitSHA,
		BuildDate: buildDate,
	}))

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
