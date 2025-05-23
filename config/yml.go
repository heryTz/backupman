package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/herytz/backupman/core"
)

type config struct {
	General struct {
		AppUrl     string `yaml:"app_url"`
		BackupCron string `yaml:"backup_cron"`
	}
	ApiKeys  []string `yaml:"api_keys"`
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
	}
	Retention struct {
		Enabled string
		By      string
		Value   int
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

func YmlToAppConfig(file string) (core.AppConfig, error) {
	c := core.AppConfig{}

	byte, err := os.ReadFile(file)
	if err != nil {
		return core.AppConfig{}, fmt.Errorf("failed to read file (%s): %s", file, err)
	}

	var ymlConfig config
	err = yaml.Unmarshal(byte, &ymlConfig)
	if err != nil {
		return core.AppConfig{}, fmt.Errorf("failed to parse yml (%s): %s", file, err)
	}

	c.General = core.GeneralConfig{
		AppUrl:     ymlConfig.General.AppUrl,
		BackupCron: ymlConfig.General.BackupCron,
	}

	c.ApiKeys = ymlConfig.ApiKeys

	switch ymlConfig.Database.Provider {
	case "mysql":
		c.Db = core.MysqlDbConfig{
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
			c.DataSources = append(c.DataSources, core.MysqlDataSourceConfig{
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
			c.Drives = append(c.Drives, core.LocalDriveConfig{
				Label:  drive.Label,
				Folder: drive.Folder,
			})
		case "google_drive":
			c.Drives = append(c.Drives, core.GoogleDriveConfig{
				Label:          drive.Label,
				Folder:         drive.Folder,
				ServiceAccount: drive.ServiceAccount,
			})
		default:
			return c, fmt.Errorf("unsupported drive provider: %s", drive.Provider)
		}
	}

	destinations := []core.MailNotifierDestinationConfig{}
	for _, dest := range ymlConfig.Notifiers.Mail.Destinations {
		destinations = append(destinations, core.MailNotifierDestinationConfig{
			Name:  dest.Name,
			Email: dest.Email,
		})
	}

	c.Notifiers = core.NotifierConfig{
		Mail: core.MailNotifierConfig{
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
			c.Retention = core.RetentionConfig{
				Enabled: true,
				Days:    ymlConfig.Retention.Value,
			}
		default:
			return c, fmt.Errorf("unsupported retention type: %s", ymlConfig.Retention.By)
		}
	}

	c.Webhooks = []core.WebhookConfig{}
	if ymlConfig.Webhooks.Enabled == "true" {
		for _, endpoint := range ymlConfig.Webhooks.Endpoints {
			c.Webhooks = append(c.Webhooks, core.WebhookConfig{
				Name:  endpoint.Name,
				Url:   endpoint.Url,
				Token: strings.TrimSpace(endpoint.Token),
			})
		}
	}

	return c, nil
}
