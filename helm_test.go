package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateDependenciesWithRefresh(t *testing.T) {
	chartPath := "charts/benchmark-subchart-level-1"
	updater := HelmUpdater{
		registryHelper: &RegistryHelper{},
	}
	startTime := time.Now()
	err := updater.UpdateChart(chartPath)
	duration := time.Since(startTime)
	assert.NoError(t, err)
	t.Logf("UpdateChart took with refresh %s", duration)
}

func TestUpdateDependenciesWithoutRefresh(t *testing.T) {
	t.Setenv("HELM_DEPS_SKIP_REFRESH", "true")
	chartPath := "charts/benchmark-subchart-level-1"
	updater := HelmUpdater{
		registryHelper: &RegistryHelper{},
	}
	updater.registryHelper.InitKubeClient()
	updater.registryHelper.UpdateRegistryInfo()
	startTime := time.Now()
	err := updater.UpdateChart(chartPath)
	duration := time.Since(startTime)
	assert.NoError(t, err)
	t.Logf("UpdateChart took with refresh %s", duration)
}
