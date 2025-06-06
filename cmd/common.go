package cmd

import (
	"github.com/herytz/backupman/cmd/config"
	"github.com/herytz/backupman/core/application"
)

func CreateAppFromYml(configFile string) (*application.App, error) {
	config, err := config.LoadYml(configFile)
	if err != nil {
		return nil, err
	}
	app := application.NewApp(config)
	return app, nil
}
