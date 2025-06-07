package cmd

import (
	"fmt"

	"github.com/herytz/backupman/core/application"
	"github.com/spf13/cobra"
)

func Version(params application.VersionConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version information",
		Long:  "Display the version information of backupman.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nCommit SHA: %s\nBuild Date: %s\n", params.Version, params.CommitSHA, params.BuildDate)
		},
	}
}
