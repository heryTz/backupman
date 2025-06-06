package http

import (
	"fmt"
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/service"
)

func SetupScheduler(app *application.App) (gocron.Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler => %s", err)
	}

	if app.Http.BackupJob.Enabled {
		job, err := scheduler.NewJob(
			gocron.CronJob(app.Http.BackupJob.Cron, true),
			gocron.NewTask(
				func(app *application.App) {
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
	}

	return scheduler, nil
}
