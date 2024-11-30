package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"helm.sh/helm/v3/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type KubeClientInterface interface {
	CoreV1() typedv1.CoreV1Interface
}

type RegistryHelper struct {
	Registries       map[string]*RegistryInfo
	Namespace        string
	kubernetesClient KubeClientInterface
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
func NewRegistryHelper(secretNames string, namespace string) *RegistryHelper {
	registryMap := make(map[string]*RegistryInfo)
	for _, registry := range strings.Split(secretNames, ",") {
		registryMap[registry] = &RegistryInfo{}
	}
	return &RegistryHelper{
		Registries: registryMap,
		Namespace:  namespace,
	}
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
	foundRegistry := r.GetRegistryByHostname(registry.Hostname)
	if foundRegistry != nil {
		action := GetRegistryAction(foundRegistry)
		log.Printf("Registry %s found in secrets, logging in", foundRegistry.SecretName)
		return action.Login()
	}
	if !registry.EnableOCI {
		log.Printf("Registry %s not found in secrets, adding it as a repo", registry.SecretName)
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
	err := r.InitKubeClient()
	if err != nil {
		return err
	}
	for secretname := range r.Registries {
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
	return nil
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
