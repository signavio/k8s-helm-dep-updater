package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func main() {
	// chartPath := flag.String("chartPath", ".", "path to the chart")
	secretNamespace := flag.String("secretNamespace", "argocd", "namespace where the secret is located")
	secretNames := flag.String("registries", os.Getenv("HELM_DEPS_SECRET_NAMES"), "comma separated list of registries to update")
	// DEPRECATED: this flag not required anymore, since the repo secret contains OCI flag already
	// if repo credentials have this flag we need to use registry login insteaad of repo add
	addRegistries := flag.String("add-registries", os.Getenv("HELM_DEPS_SECRET_NAMES_REPO_ADD"), "DEPRECATED only registries flag is required: comma separated list of registries using 'helm repo add'")
	flag.Parse()
	
	comnbinedSecretNames := strings.Join([]string{*secretNames, *addRegistries}, ",")
	registryHelper := NewRegistryHelper(comnbinedSecretNames, *secretNamespace)
	err := registryHelper.LoginAll()
	if err != nil {
		log.Fatal(err)
	}
	defer registryHelper.LogoutAll() // nolint: errcheck

	// err := updateDependencies(*chartPath, 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
