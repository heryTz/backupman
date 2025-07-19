package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/lib"
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
			fmt.Printf("Go to the following link in your browser then type the authorization code: \n\n\033[33m%s\033[0m\n\n", authURL)
			err = lib.OpenURL(authURL)
			if err != nil {
				log.Fatalf("Unable to open URL in browser: %v", err)
			}

			fmt.Print("Enter the URL you were redirected to after authorization: \n\n\033[33m")
			var redirectUrl string
			if _, err := fmt.Scan(&redirectUrl); err != nil {
				log.Fatalf("Unable to read redirect url: %v", err)
			}
			fmt.Print("\033[0m\n\n")

			u, err := url.Parse(redirectUrl)
			if err != nil {
				log.Fatalf("Unable to parse URL: %v", err)
			}

			q := u.Query()
			authCode := q.Get("code")
			if authCode == "" {
				log.Fatalf("No authorization code found in URL %s", u.String())
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			tok, err := config.Exchange(ctx, authCode)
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

	cmd.Flags().StringVar(&clientSecretFile, "client-secret-file", "google-client-secret.json", "Path to the client secret file")
	cmd.Flags().StringVar(&tokenFile, "token-file", "google-token.json", "Path that the token will be saved to")

	return cmd
}
