package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	chartPath := flag.String("chartPath", ".", "path to the chart")
	secretNamespace := flag.String("secretNamespace", "argocd", "namespace where the secret is located")
	secretNames := flag.String("registries", os.Getenv("HELM_DEPS_SECRET_NAMES"), "comma separated list of registries to update")
	addRegistries := flag.String("add-registries", os.Getenv("HELM_DEPS_SECRET_NAMES_REPO_ADD"), "comma separated list of registries using 'helm repo add'")
	ecrLoginEnabled := flag.Bool("ecr", true, "enable ecr login")
	flag.Parse()
	
	if ecrLoginEnabled != nil && *ecrLoginEnabled && *secretNames != "" {
		registryMap := make(map[string]*RegistryInfo)
		for _, registry := range strings.Split(*secretNames, ",") {
			registryMap[registry] = &RegistryInfo{}
		}
		registryHelper := &RegistryHelper{
			Registries: registryMap,
			Namespace:  *secretNamespace,
		}
		err := registryHelper.UpdateRegistryInfo()
		if err != nil {
			log.Fatal(err)
		}
		err = registryHelper.RegistryLogin()
		if err != nil {
			log.Fatal(err)
		}

		defer registryHelper.RegistryLogout() // nolint: errcheck
	}
	if addRegistries != nil && *addRegistries != "" {
		registryMap := make(map[string]*RegistryInfo)
		for _, registry := range strings.Split(*addRegistries, ",") {
			registryMap[registry] = &RegistryInfo{}
		}
		registryHelper := &RegistryHelper{
			Registries: registryMap,
			Namespace:  *secretNamespace,
		}
		err := registryHelper.UpdateRegistryInfo()
		if err != nil {
			log.Fatal(err)
		}
		err = registryHelper.RegistryAdd()
		if err != nil {
			log.Fatal(err)
		}
	}
	err := updateDependencies(*chartPath)
	if err != nil {
		log.Fatal(err)
	}
}
