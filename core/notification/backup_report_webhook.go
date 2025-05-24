package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/herytz/backupman/core"
)

func BackupReportWebhook(app *core.App, backupId string) error {
	backup, err := app.Db.Backup.ReadFullById(backupId)
	if err != nil {
		return fmt.Errorf("failed to read backup => %s", err)
	}

	for _, wh := range app.Webhooks {
		body := map[string]interface{}{
			"Event":   "backup_report",
			"Payload": backup,
		}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body => %w", err)
		}
		err = send(wh, jsonBody)
		if err != nil {
			log.Printf("failed to send webhook[backup_report] (%s) => %s", wh.Url, err)
			continue
		}
	}

	return nil
}

func send(wh core.Webhook, body []byte) error {
	client := &http.Client{}

	req, err := http.NewRequest("POST", wh.Url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request => %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Token", wh.Token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request => %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned error status: %d", resp.StatusCode)
	}

	return nil
}
