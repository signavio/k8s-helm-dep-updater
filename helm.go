package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"helm.sh/helm/v3/pkg/chart/loader"
)

type ChartLevel struct {
	Path  string
	Level int
}

func updateDependencies(chartPath string, level int) error {
	chart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}

	regex := regexp.MustCompile("file://(.*)")
	var nestedCharts []ChartLevel

	// First, collect all nested dependencies with their hierarchy level
	for _, dep := range chart.Metadata.Dependencies {
		if regex.MatchString(dep.Repository) {
			relativePath := strings.TrimPrefix(dep.Repository, "file://")
			depPath := filepath.Join(chartPath, relativePath)
			nestedCharts = append(nestedCharts, ChartLevel{Path: depPath, Level: level + 1})
		}
	}

	// Update dependencies for nested charts first (depth-first)
	var wg sync.WaitGroup
	errChan := make(chan error, len(nestedCharts))

	for _, nestedChart := range nestedCharts {
		wg.Add(1)
		go func(nestedChart ChartLevel) {
			defer wg.Done()
			fmt.Printf("Updating dependency at level %d: %s\n", nestedChart.Level, nestedChart.Path)
			if err := updateDependencies(nestedChart.Path, nestedChart.Level); err != nil {
				errChan <- err
			}
		}(nestedChart)
	}

	wg.Wait()
	close(errChan)

	// Check for any errors that occurred during the dependency updates
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	fmt.Printf("Updating chart dependencies for %s at level %d\n", chartPath, level)
	return helmDepUpdate(chartPath)
}

// Remove lock file to allow to rebuild dependencies on helm build command
func cleanupLockFiles(chartPath string) error {
	lockFiles := []string{"Chart.lock", "requirements.lock"}
	for _, lockFile := range lockFiles {
		path := filepath.Join(chartPath, lockFile)
		if _, err := os.Stat(path); err == nil {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func runHelmCommand(args ...string) error {
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func helmDepUpdate(chartPath string) error {
	err := cleanupLockFiles(chartPath)
	if err != nil {
		return err
	}
	skipRefresh := os.Getenv("HELM_DEPS_SKIP_REFRESH") // prevent update of repositories
	args := []string{"dependency", "update", chartPath}
	if skipRefresh == "true" {
		args = append(args, "--skip-refresh")
	}
	return runHelmCommand(args...)
}

