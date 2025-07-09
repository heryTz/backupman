//go:build test_integration
// +build test_integration

package service_test

import (
	"testing"

	"github.com/herytz/backupman/core/application"
	"github.com/herytz/backupman/core/lib"
	"github.com/herytz/backupman/core/service"
	"github.com/stretchr/testify/assert"
)

func TestHealthSuccess(t *testing.T) {
	app := application.NewAppMock()
	report, err := service.Health(app)
	assert.NoError(t, err)
	assert.Equal(t, lib.HEALTH_UP, report.Status)
	assert.Equal(t, lib.HEALTH_UP, report.Details["Database"].Status)
	for _, v := range app.Notifiers {
		assert.Equal(t, lib.HEALTH_UP, report.Details["Notifiers"].Components[v.GetName()].Status)
	}
	for _, v := range app.Drives {
		assert.Equal(t, lib.HEALTH_UP, report.Details["Drives"].Components[v.GetLabel()].Status)
	}
	for _, v := range app.Dumpers {
		assert.Equal(t, lib.HEALTH_UP, report.Details["DataSources"].Components[v.GetLabel()].Status)
	}
}
