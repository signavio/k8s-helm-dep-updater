package main

import (
	"fmt"
	"os/exec"
	"strings"
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

	numberObject, err := countObjects(chartPath)
	assert.NoError(t, err)
	assert.Equal(t, 40, numberObject, "Number of objects in umbrella chart")
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
	
	numberObject, err := countObjects(chartPath)
	assert.NoError(t, err)
	assert.Equal(t, 40, numberObject, "Number of objects in umbrella chart")
}

// render helm template of chart and count yaml objects
func countObjects(chartPath string) (int, error) {
	output, err := runHelmTemplate(chartPath)
	if err != nil {
		return 0, err
	}
	yamlDocs := strings.Split(output, "\n---\n")
	count := 0
	for _, doc := range yamlDocs {
		if strings.TrimSpace(doc) != "" {
			count++
		}
	}
	return count, nil
}

func runHelmTemplate(chartPath string) (string, error) {
	cmd := exec.Command("helm", "template", chartPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run helm template: %w", err)
	}
	return string(output), nil
}
