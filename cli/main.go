package main

import (
	"flag"
	"log"

	"github.com/herytz/backupman/config"
	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/service"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.yml", "Path to the config file")
	flag.Parse()

	config, err := config.YmlToAppConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	app := core.NewApp(config)
	app.Mode = core.APP_MODE_CLI

	backupIds, err := service.Backup(app)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Backups created: %v", backupIds)
}
