package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/herytz/backupman/core/application"
)

type config struct {
	Http struct {
		AppUrl    string   `yaml:"app_url"`
		ApiKeys   []string `yaml:"api_keys"`
		BackupJob struct {
			Enabled string `yaml:"enabled"`
			Cron    string `yaml:"cron"`
		} `yaml:"backup_job"`
	}
	Database struct {
		Provider string
		Host     string
		Port     int
		DbName   string `yaml:"db_name"`
		User     string
		Password string
		Tls      string
	}
	DataSources []struct {
		Provider string
		Label    string
		// mysql
		Host      string
		Port      int
		User      string
		Password  string
		DdName    string `yaml:"db_name"`
		TmpFolder string `yaml:"tmp_folder"`
		Tls       string
	} `yaml:"data_sources"`
	Drives []struct {
		Provider string
		Label    string
		// local
		Folder string
		// google drive
		ServiceAccount string `yaml:"service_account"`
	}
	Notifiers struct {
		Mail struct {
			Enabled      string
			SmtpHost     string `yaml:"smtp_host"`
			SmtpPort     int    `yaml:"smtp_port"`
			SmtpUser     string `yaml:"smtp_user"`
			SmtpPassword string `yaml:"smtp_password"`
			SmtpCrypto   string `yaml:"smtp_crypto"`
			Destinations []struct {
				Name  string
				Email string
			}
		}
		Webhooks struct {
			Enabled   string
			Endpoints []struct {
				Name  string
				Url   string
				Token string
			}
		}
	}
	Retention struct {
		Enabled string
		By      string
		Value   int
	}
}

func LoadYml(file string) (application.AppConfig, error) {
	c := application.AppConfig{}

	byte, err := os.ReadFile(file)
	if err != nil {
		return application.AppConfig{}, fmt.Errorf("failed to read file (%s): %s", file, err)
	}

	var ymlConfig config
	err = yaml.Unmarshal(byte, &ymlConfig)
	if err != nil {
		return application.AppConfig{}, fmt.Errorf("failed to parse yml (%s): %s", file, err)
	}

	httpConfig := application.HttpConfig{}
	httpConfig.AppUrl = ymlConfig.Http.AppUrl
	httpConfig.ApiKeys = ymlConfig.Http.ApiKeys
	httpConfig.BackupJob.Enabled = ymlConfig.Http.BackupJob.Enabled == "true"
	httpConfig.BackupJob.Cron = ymlConfig.Http.BackupJob.Cron
	c.Http = httpConfig

	switch ymlConfig.Database.Provider {
	case "mysql":
		c.Db = application.MysqlDbConfig{
			Host:     ymlConfig.Database.Host,
			Port:     ymlConfig.Database.Port,
			User:     ymlConfig.Database.User,
			Password: ymlConfig.Database.Password,
			Database: ymlConfig.Database.DbName,
			Tls:      ymlConfig.Database.Tls,
		}
	default:
		return c, fmt.Errorf("unsupported database provider: %s", ymlConfig.Database.Provider)
	}

	if len(ymlConfig.DataSources) == 0 {
		return c, fmt.Errorf("no data sources configured")
	}

	for _, ds := range ymlConfig.DataSources {
		switch ds.Provider {
		case "mysql":
			c.DataSources = append(c.DataSources, application.MysqlDataSourceConfig{
				Host:      ds.Host,
				Port:      ds.Port,
				User:      ds.User,
				Password:  ds.Password,
				Database:  ds.DdName,
				TmpFolder: ds.TmpFolder,
				Label:     ds.Label,
				Tls:       ds.Tls,
			})
		default:
			return c, fmt.Errorf("unsupported data source provider: %s", ds.Provider)
		}
	}

	if len(ymlConfig.Drives) == 0 {
		return c, fmt.Errorf("no drives configured")
	}

	for _, drive := range ymlConfig.Drives {
		switch drive.Provider {
		case "local":
			c.Drives = append(c.Drives, application.LocalDriveConfig{
				Label:  drive.Label,
				Folder: drive.Folder,
			})
		case "google_drive":
			c.Drives = append(c.Drives, application.GoogleDriveConfig{
				Label:          drive.Label,
				Folder:         drive.Folder,
				ServiceAccount: drive.ServiceAccount,
			})
		default:
			return c, fmt.Errorf("unsupported drive provider: %s", drive.Provider)
		}
	}

	destinations := []application.MailNotifierDestinationConfig{}
	for _, dest := range ymlConfig.Notifiers.Mail.Destinations {
		destinations = append(destinations, application.MailNotifierDestinationConfig{
			Name:  dest.Name,
			Email: dest.Email,
		})
	}

	webhooks := []application.WebhookNotifierConfig{}
	if ymlConfig.Notifiers.Webhooks.Enabled == "true" {
		for _, endpoint := range ymlConfig.Notifiers.Webhooks.Endpoints {
			webhooks = append(webhooks, application.WebhookNotifierConfig{
				Name:  endpoint.Name,
				Url:   endpoint.Url,
				Token: strings.TrimSpace(endpoint.Token),
			})
		}
	}

	c.Notifiers = application.NotifierConfig{
		Webhooks: webhooks,
		Mail: application.MailNotifierConfig{
			Enabled:      ymlConfig.Notifiers.Mail.Enabled == "true",
			SmtpHost:     ymlConfig.Notifiers.Mail.SmtpHost,
			SmtpPort:     ymlConfig.Notifiers.Mail.SmtpPort,
			SmtpUser:     ymlConfig.Notifiers.Mail.SmtpUser,
			SmtpPassword: ymlConfig.Notifiers.Mail.SmtpPassword,
			SmtpCrypto:   ymlConfig.Notifiers.Mail.SmtpCrypto,
			Destinations: destinations,
		},
	}

	if ymlConfig.Retention.Enabled == "true" {
		switch ymlConfig.Retention.By {
		case "age":
			c.Retention = application.RetentionConfig{
				Enabled: true,
				Days:    ymlConfig.Retention.Value,
			}
		default:
			return c, fmt.Errorf("unsupported retention type: %s", ymlConfig.Retention.By)
		}
	}

	return c, nil
}
