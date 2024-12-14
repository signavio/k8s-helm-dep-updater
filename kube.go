// SPDX-FileCopyrightText: 2024 Karl Theil @karlderkaefer
// SPDX-FileContributor: Karl Theil @karlderkaefer
//
// SPDX-License-Identifier: MIT-Modern-Variant

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"math/rand"

	"helm.sh/helm/v3/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	ArgocdRepoSecretLabel = "argocd.argoproj.io/secret-type"
	ArgocdRepoSecretValue = "repository"
)

type KubeClientInterface interface {
	CoreV1() typedv1.CoreV1Interface
}

type RegistryHelper struct {
	Registries       map[string]*RegistryInfo
	Namespace        string
	kubernetesClient KubeClientInterface
	config           *HelmUpdateConfig
}

type RegistryInfo struct {
	Hostname   string
	Username   string
	Password   string
	SecretName string
	EnableOCI  bool
}

// RegistryAction interface for registry operations
type RegistryAction interface {
	Login() error
	Logout() error
}

// OCIStrategy for OCI-enabled registries
type OCIStrategy struct {
	*RegistryInfo
}

func (o *OCIStrategy) Login() error {
	return runHelmCommand("registry", "login", "-u", o.Username, "-p", o.Password, o.Hostname)
}

func (o *OCIStrategy) Logout() error {
	return runHelmCommand("registry", "logout", o.Hostname)
}

// DefaultStrategy for non-OCI registries
type DefaultStrategy struct {
	*RegistryInfo
}

func (d *DefaultStrategy) Login() error {
	return runHelmCommand("repo", "add", d.SecretName, d.Hostname, "--username", d.Username, "--password", d.Password)
}

func (d *DefaultStrategy) Logout() error {
	// Logout is not needed for DefaultStrategy
	return nil
}

// NewRegistryHelper creates a new RegistryHelper
// secretNames is a comma separated list of registry secrets
func NewRegistryHelper(secretNames string, namespace string, config *HelmUpdateConfig) *RegistryHelper {
	registryMap := make(map[string]*RegistryInfo)
	for _, registry := range strings.Split(secretNames, ",") {
		if registry != "" {
			registryMap[registry] = &RegistryInfo{}
		}
	}
	r := &RegistryHelper{
		Registries: registryMap,
		Namespace:  namespace,
		config:     config,
	}
	if config.UseRandomHelmCacheDir {
		r.SetRandomHelmCacheDir()
	}
	return r
}

func (r *RegistryHelper) SetRandomHelmCacheDir() error {
	// Create a local rand generator
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Create unique temporary directories
	cacheDir := fmt.Sprintf("/tmp/helm_cache_%d", seed.Int())
	configDir := fmt.Sprintf("/tmp/helm_config_%d", seed.Int())

	// Ensure the directories exist
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	// Set the environment variables for Helm
	os.Setenv("HELM_REPOSITORY_CACHE", cacheDir)
	os.Setenv("HELM_REPOSITORY_CONFIG", filepath.Join(configDir, "repositories.yaml"))
	return nil
}

// GetRegistry returns the credential protected registry by hostname
func (r *RegistryHelper) GetRegistryByHostname(registry string) *RegistryInfo {
	for _, registryInfo := range r.Registries {
		if registryInfo.Hostname == registry {
			return registryInfo
		}
	}
	return nil
}

// LoginIfExists logs into the registry if the hostname can found in existing secrets
func (r *RegistryHelper) LoginIfExists(registry *RegistryInfo) error {
	if registry == nil {
		return errors.New("registry can not be empty")
	}
	exists, err := helmRepoExists(registry, r.config)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Registry %s found in helm repos, skipping login", registry.SecretName)
		return nil
	}
	foundRegistry := r.GetRegistryByHostname(registry.Hostname)
	if foundRegistry != nil {
		action := GetRegistryAction(foundRegistry)
		log.Printf("Registry %s found in secrets, logging in", foundRegistry.SecretName)
		return action.Login()
	}
	if !registry.EnableOCI {
		log.Printf("Registry %s ≈", registry.SecretName)
		action := GetRegistryAction(registry)
		return action.Login()
	}
	return nil
}

func (r *RegistryHelper) InitKubeClient() error {
	kubeConfig := kube.GetConfig("", "", "")
	restConfig, err := kubeConfig.ToRESTConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	r.kubernetesClient = clientset
	return nil
}

func (r *RegistryHelper) UpdateRegistryInfo() error {
	if len(r.Registries) == 0 && !r.config.FetchArgocdRepoSecrets {
		log.Printf("No secrets provided, skipping registry update from kubeclient")
		return nil
	}
	for secretname := range r.Registries {
		err := r.InitKubeClient()
		if err != nil {
			return err
		}
		secret, err := r.kubernetesClient.CoreV1().Secrets(r.Namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
		if err != nil {
			return err
		}
		r.Registries[secretname] = &RegistryInfo{
			Hostname:   string(secret.Data["url"]),
			Username:   string(secret.Data["username"]),
			Password:   string(secret.Data["password"]),
			EnableOCI:  string(secret.Data["enableOCI"]) == "true",
			SecretName: secretname,
		}
	}
	if r.config.FetchArgocdRepoSecrets {
		err := r.InitKubeClient()
		if err != nil {
			return err
		}
		r.SetRegistriesByLabel()
	}
	log.Printf("Created %d registries", len(r.Registries))
	return nil
}

func (r *RegistryHelper) SetRegistriesByLabel() {
	selector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			ArgocdRepoSecretLabel: ArgocdRepoSecretValue,
		},
	}
	secrets, err := r.kubernetesClient.CoreV1().Secrets(r.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&selector),
	})
	if err != nil {
		log.Printf("Unable to list secrets with label %s: %v", ArgocdRepoSecretLabel, err)
		return
	}
	for _, secret := range secrets.Items {
		secretName := secret.Name
		r.Registries[secretName] = &RegistryInfo{
			Hostname:   string(secret.Data["url"]),
			Username:   string(secret.Data["username"]),
			Password:   string(secret.Data["password"]),
			EnableOCI:  string(secret.Data["enableOCI"]) == "true",
			SecretName: secretName,
		}
	}
}

func GetRegistryAction(registry *RegistryInfo) RegistryAction {
	if registry.EnableOCI {
		return &OCIStrategy{RegistryInfo: registry}
	}
	return &DefaultStrategy{RegistryInfo: registry}
}

func (r *RegistryHelper) LoginAll() error {
	for _, registry := range r.Registries {
		action := GetRegistryAction(registry)
		err := action.Login()
		if err != nil {
			log.Printf("Unable to login to registry %s: %v", registry.SecretName, err)
			return err
		}
	}
	return nil
}

func (r *RegistryHelper) LogoutAll() error {
	for _, registry := range r.Registries {
		action := GetRegistryAction(registry)
		err := action.Logout()
		if err != nil {
			return err
		}
	}
	return nil
}
