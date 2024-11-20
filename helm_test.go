package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateDependencies(t *testing.T) {
	chartPath := "charts/benchmark-subchart-level-1"
	timedUpdateDependencies := measureTime(updateDependencies)

	os.Setenv("HELM_DEPS_SKIP_REFRESH", "true")
	duration, err := timedUpdateDependencies(chartPath, 1)
	assert.NoError(t, err)
	t.Logf("updateDependencies without refresh took %s", duration)
	cleanupLockFiles(chartPath)
	
	os.Setenv("HELM_DEPS_SKIP_REFRESH", "false")
	duration, err = timedUpdateDependencies(chartPath, 1)
	assert.NoError(t, err)
	t.Logf("updateDependencies with refresh took %s", duration)
	cleanupLockFiles(chartPath)


}

// Decorator function to measure execution time
func measureTime(fn func(string, int) error) func(string, int) (time.Duration, error) {
	return func(chartPath string, level int) (time.Duration, error) {
		start := time.Now()
		err := fn(chartPath, level)
		duration := time.Since(start)
		// cleanupLockFiles(chartPath)
		return duration, err
	}
}
