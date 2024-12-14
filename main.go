// SPDX-FileCopyrightText: 2024 Karl Theil @karlderkaefer
// SPDX-FileContributor: Karl Theil @karlderkaefer
//
// SPDX-License-Identifier: MIT-Modern-Variant

package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	chartPath := flag.String("chartPath", ".", "path to the chart")
	secretNamespace := flag.String("secretNamespace", "argocd", "namespace where the secret is located")
	secretNames := flag.String("registries", os.Getenv("HELM_DEPS_SECRET_NAMES"), "comma separated list of registries to update")
	// DEPRECATED: this flag not required anymore, since the repo secret contains OCI flag already
	// if repo credentials have this flag we need to use registry login insteaad of repo add
	addRegistries := flag.String("add-registries", os.Getenv("HELM_DEPS_SECRET_NAMES_REPO_ADD"), "DEPRECATED only registries flag is required: comma separated list of registries using 'helm repo add'")

	// this will skip adding the repo if it already exists. useful during testing, if the depdenccy are already added
	skipRepoOverwrite := parseBoolEnv("HELM_DEPS_SKIP_REPO_OVERWRITE", false)
	// skips repo update in sub dependencies and give massive performance improvement
	skipDepRefresh := parseBoolEnv("HELM_DEPS_SKIP_REFRESH", false)
	// you can skip the login to all available registry at the start in combination with HELM_DEPS_SKIP_REFRESH=true
	skipLoginAtStart := parseBoolEnv("HELM_DEPS_SKIP_START_LOGIN", false)
	skipDepRefreshFlag := flag.Bool("skip-dep-refresh", skipDepRefresh, "Set environment variable HELM_DEPS_SKIP_REFRESH to use this flag (default false). run repo refresh in every sub dependency. this will add the helm argument --skip-refresh to helm dependency update")
	skipRepoOverwriteFlag := flag.Bool("skip-repo-overwrite", skipRepoOverwrite, "Env: HELM_DEPS_SKIP_REPO_OVERWRITE (default true). This will not add a repo if it already exists")
	skipLoginAtStartFlag := flag.Bool("skip-login-at-start", skipLoginAtStart, "Env: HELM_DEPS_SKIP_START_LOGIN (default false). This will skip login to all available registry at the start in combination with HELM_DEPS_SKIP_REFRESH=true")
	fetchArgocdRepoSecrets := parseBoolEnv("HELM_DEPS_FETCH_ARGOCD_REPO_SECRETS", false)
	fetchArgocdRepoSecretsFlag := flag.Bool("fetch-argocd-repo-secrets", fetchArgocdRepoSecrets, "Env: HELM_DEPS_FETCH_ARGOCD_REPO_SECRETS (default false). Fetch the argocd repository secrets as registries")
	useRandomHelmCacheDir := parseBoolEnv("HELM_DEPS_RANDOM_CACHE_DIR", false)
	useRandomHelmCacheDirFlag := flag.Bool("use-random-helm-cache-dir", useRandomHelmCacheDir, "Env: HELM_DEPS_RANDOM_CACHE_DIR (default false). Use a random cache directory for helm")

	flag.Parse()

	config := &HelmUpdateConfig{
		FetchArgocdRepoSecrets: *fetchArgocdRepoSecretsFlag,
		SkipRepoOverwrite:      *skipRepoOverwriteFlag,
		SkipDepdencyRefresh:    *skipDepRefreshFlag,
		UseRandomHelmCacheDir:  *useRandomHelmCacheDirFlag,
	}
	// just for backward compatibility
	comnbinedSecretNames := strings.Join([]string{*secretNames, *addRegistries}, ",")
	registryHelper := NewRegistryHelper(comnbinedSecretNames, *secretNamespace, config)

	err := registryHelper.UpdateRegistryInfo()
	if err != nil {
		log.Fatal("Unable to update registry info: ", err)
	}

	if !*skipLoginAtStartFlag {
		err := registryHelper.LoginAll()
		if err != nil {
			log.Fatal("Unable to Login: ", err)
		}
		defer registryHelper.LogoutAll() // nolint: errcheck
	}
	updater := HelmUpdater{
		registryHelper: registryHelper,
		config:         config,
	}
	err = updater.UpdateChart(*chartPath)
	if err != nil {
		log.Fatal("Could not update the chart", err)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// parseBoolEnv parses an environment variable as a boolean, with a default fallback.
func parseBoolEnv(envVar string, defaultValue bool) bool {
	val := getEnv(envVar, strconv.FormatBool(defaultValue))
	parsedVal, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}
	return parsedVal
}
