package cmd

import (
	"github.com/herytz/backupman/config"
	"github.com/herytz/backupman/core"
)

func CreateAppFromYml(configFile string) (*core.App, error) {
	config, err := config.YmlToAppConfig(configFile)
	if err != nil {
		return nil, err
	}
	app := core.NewApp(config)
	return app, nil
}
