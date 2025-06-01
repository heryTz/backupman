package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type VersionParams struct {
	Version   string
	CommitSHA string
	BuildDate string
}

func Version(params VersionParams) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version information",
		Long:  "Display the version information of backupman.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nCommit SHA: %s\nBuild Date: %s\n", params.Version, params.CommitSHA, params.BuildDate)
		},
	}
}
