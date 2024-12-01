package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"helm.sh/helm/v3/pkg/chart/loader"
)

type ChartInfo struct {
	Path         string
	Level        int
	NestedCharts []ChartInfo
	Registries   map[string]*RegistryInfo
}
type HelmUpdater struct {
	registryHelper *RegistryHelper
	config *HelmUpdateConfig
}

type HelmUpdateConfig struct {
	SkipRepoOverwrite bool
	SkipDepdencyRefresh bool // ENV: HELM_DEPS_SKIP_REFRESH
}

func (c *ChartInfo) AddDependencyUrl(depdencyUrl string) error {
	if depdencyUrl == "" {
		return nil
	}
	hostnameKey := sanitizeHostname(depdencyUrl)
	parsedURL, err := url.Parse(depdencyUrl)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %s in Chart %s : %w", depdencyUrl, c.Path, err)
	}
	if parsedURL.Scheme == "file" {
		return nil
	}
	if !(parsedURL.Scheme == "oci" || parsedURL.Scheme == "https") {
		log.Printf("Unsupported scheme %s in URL: %s in Chart %s. Skipping adding registry", parsedURL.Scheme, depdencyUrl, c.Path)
		return nil
	}
	registry := &RegistryInfo{
		Hostname:   parsedURL.String(),
		SecretName: hostnameKey,
		EnableOCI:  false,
	}
	if parsedURL.Scheme == "oci" {
		registry.EnableOCI = true
	}
	c.Registries[hostnameKey] = registry
	return nil
}

func (updater *HelmUpdater) UpdateChart(chartPath string) error {
	chartInfo, err := loadChartInfo(chartPath, 1)
	if err != nil {
		return err
	}
	// we login to all collected registries before updating dependencies
	if updater.config.SkipDepdencyRefresh {
		for _, registry := range chartInfo.Registries {
			err := updater.registryHelper.LoginIfExists(registry)
			if err != nil {
				return err
			}
		}
		// update helm repo after adding all registries
		err := runHelmCommand("repo", "update")
		if err != nil {
			return err
		}
	}
	return updater.updateDependencies(chartInfo)
}

// sanitizeHostname replaces all non-alphanumeric characters with hyphens
func sanitizeHostname(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return re.ReplaceAllString(input, "-")
}

func loadChartInfo(chartPath string, level int) (ChartInfo, error) {
	chart, err := loader.Load(chartPath)
	if err != nil {
		return ChartInfo{}, err
	}

	chartInfo := ChartInfo{Path: chartPath, Level: level, Registries: make(map[string]*RegistryInfo)}
	regex := regexp.MustCompile("file://(.*)")

	// Collect URLs of dependencies
	for _, dep := range chart.Metadata.Dependencies {
		chartInfo.AddDependencyUrl(dep.Repository)
		if regex.MatchString(dep.Repository) {
			relativePath := strings.TrimPrefix(dep.Repository, "file://")
			depPath := filepath.Join(chartPath, relativePath)
			nestedChart, err := loadChartInfo(depPath, level+1)
			if err != nil {
				return ChartInfo{}, err
			}
			chartInfo.NestedCharts = append(chartInfo.NestedCharts, nestedChart)
			// Merge URLs from nested charts
			for _, registry := range nestedChart.Registries {
				chartInfo.AddDependencyUrl(registry.Hostname)
			}
		}
	}

	return chartInfo, nil
}

func (updater *HelmUpdater) updateDependencies(chartInfo ChartInfo) error {
	// Update dependencies for nested charts
	var wg sync.WaitGroup
	errChan := make(chan error, len(chartInfo.NestedCharts))

	for _, nestedChart := range chartInfo.NestedCharts {
		wg.Add(1)
		go func(nestedChart ChartInfo) {
			defer wg.Done()
			fmt.Printf("Updating chart dependencies for %s at level %d\n", nestedChart.Path, nestedChart.Level)
			if err := updater.updateDependencies(nestedChart); err != nil {
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
	return updater.helmDepUpdate(chartInfo.Path)
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

// helmRepoExists checks if a helm repository already exists with helm repo ls command
func helmRepoExists(registry *RegistryInfo, config *HelmUpdateConfig) (bool, error) {
	if !config.SkipRepoOverwrite {
		return false, nil
	}
	cmd := exec.Command("helm", "repo", "ls")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("failed to run helm repo ls: %w", err)
	}
	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, registry.Hostname) {
			return true, nil
		}
	}
	return false, nil
}

func (updater *HelmUpdater) helmDepUpdate(chartPath string) error {
	err := cleanupLockFiles(chartPath)
	if err != nil {
		return err
	}
	args := []string{"dependency", "update", chartPath}
	if updater.config.SkipDepdencyRefresh {
		args = append(args, "--skip-refresh")
	}
	return runHelmCommand(args...)
}
