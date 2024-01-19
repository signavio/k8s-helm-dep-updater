# Helm Dependency Updater

This project resolves the problem of missing recursive dependency updates in Helm charts. It is a simple tool that can be used to update the dependencies of a Helm chart recursively.

There is an [open issue](https://github.com/helm/helm/issues/2247) and [PR](https://github.com/helm/helm/pull/11766) in the Helm project to add this functionality to Helm itself. Until this is resolved, this tool can be used to update the dependencies of a Helm chart recursively. It will be part of the [Milestone 3.13.0](https://github.com/helm/helm/milestone/131) which is released on 13th of September 2023.

## Install

```bash
wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/download/v1.1.1/k8s-helm-dep-updater_$(uname -s)_$(uname -m).tar.gz" | tar -C /tmp -xzf- k8s-helm-dep-updater
sudo mv /tmp/k8s-helm-dep-updater /usr/local/bin/k8s-helm-dep-updater
```

Install Helm Plugin

```bash
helm plugin install https://github.com/signavio/k8s-helm-dep-updater.git
```

## Argocd Integration

To activate the helm dep updater we use downloader plugin to activate it. `dummy` can be anything. Most importantly with the `deps://` prefix we activate the downloader plugin.

```yaml
project: default
source:
  repoURL: 'https://github.com/signavio/k8s-helm-dep-updater'
  path: charts/umbrella
  targetRevision: main
  helm:
    valueFiles:
      - 'deps://dummy'
destination:
  server: 'https://kubernetes.default.svc'
  namespace: default
```

To add the plugin to argocd, those are the relevant lines:

```yaml
repoServer:
  initContainers:
    - name: download-tools
      image: alpine:3
      command: [sh, -c]
      args:
        - |
          mkdir -p /custom-tools/helm-plugins/helm-dep-updater
          wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/download/v1.1.1/k8s-helm-dep-updater.tar.gz" | tar -C /custom-tools/helm-plugins/helm-dep-updater -xzf-;
          wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/download/v1.1.1/k8s-helm-dep-updater_$(uname -s)_$(uname -m).tar.gz" | tar -C /custom-tools/ -xzf- k8s-helm-dep-updater
    volumeMounts:
      - mountPath: /custom-tools
        name: custom-tools
  env:
  - name: HELM_PLUGINS
    value: /custom-tools/helm-plugins
  # optionally if you want to login to ECR before fetching the charts
  # it expects a comma separated list of secret names
  - name: HELM_DEPS_SECRET_NAMES
    value: helm-ecr-staging,helm-ecr-production
  volumeMounts:
    - mountPath: /custom-tools
      name: custom-tools
    - mountPath: /usr/local/bin/k8s-helm-dep-updater
      name: custom-tools
      subPath: k8s-helm-dep-updater
```

## Usage

```bash
helm template . -f deps://dummy
```
