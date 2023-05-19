package main

import (
	"context"

	"helm.sh/helm/v3/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RegistryHelper struct {
	Registries map[string]*RegistryInfo
	Namespace  string
}

type RegistryInfo struct {
	Hostname string
	Username string
	Password string
}

func (r *RegistryInfo) Login() error {
	return runHelmCommand("registry", "login", "-u", r.Username, "-p", r.Password, r.Hostname)
}

func (r *RegistryInfo) Logout() error {
	return runHelmCommand("registry", "logout", r.Hostname)
}

func NewKubeClient() (*kubernetes.Clientset, error) {
	kubeConfig := kube.GetConfig("", "", "")
	restConfig, err := kubeConfig.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func (r *RegistryHelper) UpdateRegistryInfo() error {
	clientset, err := NewKubeClient()
	if err != nil {
		return err
	}
	for secretname, _ := range r.Registries {
		secret, err := clientset.CoreV1().Secrets(r.Namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
		if err != nil {
			return err
		}
		r.Registries[secretname] = &RegistryInfo{
			Hostname: string(secret.Data["url"]),
			Username: string(secret.Data["username"]),
			Password: string(secret.Data["password"]),
		}
	}
	return nil
}

func (r *RegistryHelper) RegistryLogin() error {
	for _, registry := range r.Registries {
		registry.Login()
	}
	return nil
}

func (r *RegistryHelper) RegistryLogout() error {
	for _, registry := range r.Registries {
		registry.Logout()
	}
	return nil
}
