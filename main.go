package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	chartPath := flag.String("chartPath", ".", "path to the chart")
	secretNamespace := flag.String("secretNamespace", "argocd", "namespace where the secret is located")
	secretNames := flag.String("registries", "helm-ecr-staging", "comma separated list of registries to update")
	ecrLoginEnabled := flag.Bool("ecr", false, "enable ecr login")
	flag.Parse()

	if ecrLoginEnabled != nil && *ecrLoginEnabled {
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
	err := updateDependencies(*chartPath)
	if err != nil {
		log.Fatal(err)
	}
}
