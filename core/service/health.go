package service

import (
	"log"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/lib"
)

type ComponentStatus struct {
	Status     string
	Components map[string]ComponentStatus
}

type HealthReport struct {
	Version   string
	CommitSHA string
	BuildDate string
	Status    string
	Details   map[string]ComponentStatus
}

func Health(app *application.App) (HealthReport, error) {
	report := HealthReport{
		Status:    lib.HEALTH_UP,
		Version:   app.Version.Version,
		CommitSHA: app.Version.CommitSHA,
		BuildDate: app.Version.BuildDate,
	}
	details := make(map[string]ComponentStatus)

	databaseStatus := lib.HEALTH_UP
	err := app.Db.Health.Check()
	if err != nil {
		log.Printf("Database health check failed => %s", err)
		databaseStatus = lib.HEALTH_DOWN
		report.Status = lib.HEALTH_DOWN
	}
	details["Database"] = ComponentStatus{
		Status: databaseStatus,
	}

	if app.Notifiers.Mail != nil {
		mailStatus := lib.HEALTH_UP
		err = app.Notifiers.Mail.Health()
		if err != nil {
			log.Printf("Mail health check failed => %s", err)
			mailStatus = lib.HEALTH_DOWN
		}
		details["Mail"] = ComponentStatus{
			Status: mailStatus,
		}
	}

	driveComponents := make(map[string]ComponentStatus)
	driveStatus := lib.HEALTH_UP
	for _, drive := range app.Drives {
		driveItemStatus := lib.HEALTH_UP
		err = drive.Health()
		if err != nil {
			log.Printf("Drive health check failed for %s => %s", drive.GetLabel(), err)
			driveItemStatus = lib.HEALTH_DOWN
			driveStatus = lib.HEALTH_DOWN
		}
		driveComponents[drive.GetLabel()] = ComponentStatus{
			Status: driveItemStatus,
		}
	}
	details["Drives"] = ComponentStatus{
		Status:     driveStatus,
		Components: driveComponents,
	}

	dataSourceComponents := make(map[string]ComponentStatus)
	dataSourceStatus := lib.HEALTH_UP
	for _, dataSource := range app.Dumpers {
		dataSourceItemStatus := lib.HEALTH_UP
		err = dataSource.Health()
		if err != nil {
			log.Printf("Data source health check failed for %s => %s", dataSource.GetLabel(), err)
			dataSourceItemStatus = lib.HEALTH_DOWN
			dataSourceStatus = lib.HEALTH_DOWN
		}
		dataSourceComponents[dataSource.GetLabel()] = ComponentStatus{
			Status: dataSourceItemStatus,
		}
	}
	details["DataSources"] = ComponentStatus{
		Status:     dataSourceStatus,
		Components: dataSourceComponents,
	}

	for _, d := range details {
		if d.Status == lib.HEALTH_DOWN {
			report.Status = lib.HEALTH_DOWN
			break
		}
	}

	report.Details = details
	return report, nil
}
