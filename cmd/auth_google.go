package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/herytz/backupman/core/application"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func AuthGoogle(versionConfig application.VersionConfig) *cobra.Command {
	var clientSecretFile string
	var tokenFile string

	cmd := &cobra.Command{
		Use:   "auth-google",
		Short: "Authenticate with Google Drive using OAuth2",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := os.ReadFile(clientSecretFile)
			if err != nil {
				log.Fatalf("Unable to read client secret file: %v", err)
			}

			config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
			if err != nil {
				log.Fatalf("Unable to parse client secret file to config: %v", err)
			}

			authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			fmt.Printf("Go to the following link in your browser then type the "+
				"authorization code: \n%v\n", authURL)

			var authCode string
			if _, err := fmt.Scan(&authCode); err != nil {
				log.Fatalf("Unable to read authorization code: %v", err)
			}

			tok, err := config.Exchange(context.TODO(), authCode)
			if err != nil {
				log.Fatalf("Unable to retrieve token from web: %v", err)
			}

			fmt.Printf("Saving credential file to: %s\n", tokenFile)
			f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Fatalf("Unable to cache oauth token: %v", err)
			}
			defer f.Close()
			json.NewEncoder(f).Encode(tok)
		},
	}

	cmd.Flags().StringVar(&clientSecretFile, "client-secret-file", "client-secret-file", "Path to the client secret file")
	cmd.Flags().StringVar(&tokenFile, "token-file", "token-file", "Path to the token file")
	// cmd.MarkFlagRequired("client-secret-file")
	// cmd.MarkFlagRequired("token-file")

	return cmd
}
