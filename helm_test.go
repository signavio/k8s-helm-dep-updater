package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	name                string
	chartPath           string
	helmDepsSkipRefresh bool
	expectedNumObjects  int
}

func TestUpdateDependencies(t *testing.T) {
	testCases := []TestCase{
		{
			name:                "With Refresh",
			chartPath:           "charts/benchmark-subchart-level-1",
			helmDepsSkipRefresh: false,
			expectedNumObjects:  42,
		},
		{
			name:                "Without Refresh",
			chartPath:           "charts/benchmark-subchart-level-1",
			helmDepsSkipRefresh: true,
			expectedNumObjects:  42,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &HelmUpdateConfig{}
			if tc.helmDepsSkipRefresh {
				config.SkipDepdencyRefresh = true
				config.SkipRepoOverwrite = true
			}

			updater := HelmUpdater{
				registryHelper: &RegistryHelper{
					config: config,
				},
				config: config,
			}

			if tc.helmDepsSkipRefresh {
				updater.registryHelper.UpdateRegistryInfo()
			}

			startTime := time.Now()
			err := updater.UpdateChart(tc.chartPath)
			duration := time.Since(startTime)
			assert.NoError(t, err)
			t.Logf("UpdateChart took %s", duration)

			numberObject, err := countObjects(tc.chartPath)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedNumObjects, numberObject, "Number of objects in umbrella chart")
		})
	}
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
