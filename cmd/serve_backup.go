package cmd

import (
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/http"
	"github.com/spf13/cobra"
)

func ServeBackup(version application.VersionConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the backup manager",
		Long:  "This command will start the backup manager server.",
		Run: func(cmd *cobra.Command, args []string) {
			configFile, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal(err)
			}
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				log.Fatal(err)
			}

			app, err := CreateAppFromYml(configFile)
			if err != nil {
				log.Fatalf("Error creating app from config => %v", err)
			}
			app.Mode = application.APP_MODE_WEB
			app.Version = version
			err = http.Serve(app, port)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().IntP("port", "p", 8080, "Port to run the server on")

	return cmd
}
