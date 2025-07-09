package message

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/herytz/backupman/core/model"
)

const tmpl = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        .container {
            max-width: 800px;
            margin: 20px auto;
            font-family: Arial, sans-serif;
        }
        .header {
            border-bottom: 2px solid #007bff;
            padding-bottom: 10px;
            margin-bottom: 25px;
        }
        .status-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        .status-table td {
            padding: 12px;
            border: 1px solid #ddd;
        }
        .status-header {
            background-color: #f8f9fa;
            font-weight: bold;
        }
        .status-badge {
            padding: 5px 10px;
            border-radius: 12px;
            font-size: 0.9em;
        }
        .success {
            background-color: #d4edda;
            color: #155724;
        }
        .failed {
            background-color: #f8d7da;
            color: #721c24;
        }
        .pending {
            background-color: #fff3cd;
            color: #856404;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2 style="color: #007bff;">📦 Backup Report Notification</h2>
            <p style="color: #6c757d;">Your database backup report details</p>
        </div>

        <!-- Backup Metadata -->
        <table style="width: 100%; margin-bottom: 25px;">
            <tr>
                <td style="width: 30%; padding: 8px; background-color: #f8f9fa;">Backup ID</td>
                <td style="padding: 8px;">{{.BackupID}}</td>
            </tr>
            <tr>
                <td style="padding: 8px; background-color: #f8f9fa;">Backup Date</td>
                <td style="padding: 8px;">{{.BackupDate}}</td>
            </tr>
            <tr>
                <td style="padding: 8px; background-color: #f8f9fa;">Database Name</td>
                <td style="padding: 8px;">{{.DatabaseName}}</td>
            </tr>
        </table>

        <!-- Upload Status -->
        <h4 style="color: #007bff; margin-bottom: 15px;">Storage Providers Status</h4>
        <table class="status-table">
            <tr class="status-header">
                <td>Provider</td>
                <td>Status</td>
            </tr>
            {{range .UploadStatus}}
            <tr>
                <td>{{.Provider}}</td>
                <td>
                    <span class="status-badge {{.Status | ToLower}}">
                        {{.Status}}
                    </span>
                </td>
            </tr>
            {{end}}
        </table>

        <!-- Footer -->
        <div style="margin-top: 30px; color: #6c757d; font-size: 0.9em;">
            <hr style="border-top: 1px solid #eee;">
            <p>This is an automated notification. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`

type UploadStatus struct {
	Provider string
	Status   string
}

type EmailData struct {
	BackupID     string
	BackupDate   string
	DatabaseName string
	UploadStatus []UploadStatus
}

func BackupReportMail(backup *model.BackupFull) (string, error) {
	tm, err := template.
		New("backup_report").
		Funcs(template.FuncMap{"ToLower": strings.ToLower}).
		Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template => %s", err)
	}

	data := EmailData{
		BackupID:     backup.Id,
		BackupDate:   backup.CreatedAt.Format("2006-01-02 15:04:05"),
		DatabaseName: backup.Label,
	}
	for _, driveFile := range backup.DriveFiles {
		data.UploadStatus = append(data.UploadStatus, UploadStatus{
			Provider: driveFile.Provider,
			Status:   driveFile.Status,
		})
	}

	var buf bytes.Buffer
	err = tm.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %s", err)
	}

	return buf.String(), nil
}
