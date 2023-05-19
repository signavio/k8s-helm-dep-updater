package main

import (
	"context"

	"helm.sh/helm/v3/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type KubeClientInterface interface {
	CoreV1() typedv1.CoreV1Interface
}

type RegistryHelper struct {
	Registries map[string]*RegistryInfo
	Namespace  string
	kubernetesClient KubeClientInterface
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
	for secretname, _ := range r.Registries {
		secret, err := r.kubernetesClient.CoreV1().Secrets(r.Namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
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
