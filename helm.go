package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"helm.sh/helm/v3/pkg/chart/loader"
)

func updateDependencies(chartPath string) error {
	chart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}
	for _, dep := range chart.Metadata.Dependencies {
		fmt.Printf("detected dependency: %s repository: %s version: %s\n", dep.Name, dep.Repository, dep.Version)
		regex := regexp.MustCompile("file://(.*)")
		if regex.MatchString(dep.Repository) {
			relativePath := strings.TrimPrefix(dep.Repository, "file://")
			depPath := filepath.Join(chartPath, relativePath)
			fmt.Println("updating dependency:", dep.Name)
			if err := updateDependencies(depPath); err != nil {
				return err
			}
		}
	}

	return helmDepUpdate(chartPath)
}

func runHelmCommand(args ...string) error {
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func helmDepUpdate(chartPath string) error {
	return runHelmCommand("dependency", "update", chartPath)
}
