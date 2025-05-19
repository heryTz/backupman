package main

import (
	"fmt"
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/herytz/backupman/core"
	"github.com/herytz/backupman/core/service"
)

func SetupScheduler(app *core.App) (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler => %s", err)
	}
	job, err := scheduler.NewJob(
		gocron.CronJob(app.Config.BackupCron, true),
		gocron.NewTask(
			func(app *core.App) {
				log.Println("running scheduled backup...")
				backupIds, err := service.Backup(app)
				if err != nil {
					log.Printf("%s", err)
				} else {
					log.Printf("scheduled backup started, %v backups created", backupIds)
				}
			},
			app,
		),
	)
	if err != nil {
		return scheduler, fmt.Errorf("Failed to create job => %s", err)
	}

	log.Printf("scheduler job created, ID=%s\n", job.ID())
	return scheduler, nil
}
