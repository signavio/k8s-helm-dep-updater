# Helm Dependency Updater

This project resolves the problem of missing recursive dependency updates in Helm charts. It is a simple tool that can be used to update the dependencies of a Helm chart recursively. Charts from Registry already have their subcharts included, but if you use local file references `file://` you need to update subcharts recursivly if you have more than one hierarchy level of dependencies.

There is an [open issue](https://github.com/helm/helm/issues/2247) and [PR](https://github.com/helm/helm/pull/11766) in the Helm project to add this functionality to Helm itself. Until this is resolved, this tool can be used to update the dependencies of a Helm chart recursively. It will be part of the [Milestone 3.13.0](https://github.com/helm/helm/milestone/131) which is released on 13th of September 2023.

## Install

```bash
wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/latest/download/k8s-helm-dep-updater_$(uname -s)_$(uname -m).tar.gz" | tar -C /tmp -xzf- k8s-helm-dep-updater
sudo mv /tmp/k8s-helm-dep-updater /usr/local/bin/k8s-helm-dep-updater
```

Install Helm Plugin

```bash
helm plugin install https://github.com/signavio/k8s-helm-dep-updater.git
```

## Argocd Integration

To activate the helm dep updater we use downloader plugin to activate it. `dummy` can be anything. Most importantly with the `deps://` prefix you can activate the downloader plugin.

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
          wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/latest/download/k8s-helm-dep-updater.tar.gz" | tar -C /custom-tools/helm-plugins/helm-dep-updater -xzf-;
          wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/latest/download/k8s-helm-dep-updater_$(uname -s)_$(uname -m).tar.gz" | tar -C /custom-tools/ -xzf- k8s-helm-dep-updater
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

## Environment Variables

| Environment Name              | Default Value       | Description                                                                                              |
|-------------------------------|---------------------|----------------------------------------------------------------------------------------------------------|
| HELM_DEPS_SECRET_NAMES        | ""                  | Comma-separated list of registries to update.                                                            |
| HELM_DEPS_SECRET_NAMES_REPO_ADD | ""                | DEPRECATED: Comma-separated list of registries using 'helm repo add'.                                    |
| HELM_DEPS_SKIP_REPO_OVERWRITE | false               | Skips adding the repository if it already exists, useful during testing.                                 |
| HELM_DEPS_SKIP_REFRESH        | false               | Skips repository update for sub-dependencies, providing a significant performance boost.                 |
| HELM_DEPS_SKIP_START_LOGIN    | false               | Skips login to all available registries at the start if combined with HELM_DEPS_SKIP_REFRESH=true.       |
| HELM_DEPS_FETCH_ARGOCD_REPO_SECRETS | false | Adds registries that are registered in the argocd repository secrets.                                     |

### Local

Local Usage:

```bash
export HELM_DEPS_SKIP_REPO_OVERWRITE=true
cd charts/benchmark-subchart-level-1
k8s-helm-dep-updater

# for maximum performance
export HELM_DEPS_SKIP_REFRESH=true
export HELM_DEPS_SKIP_REPO_OVERWRITE=true
export HELM_DEPS_SKIP_START_LOGIN=true
k8s-helm-dep-updater
```

### Helm Plugin Usage

Usage as helm plugin is the same. All Configuration needs to be done via environment variables.

```bash
export HELM_DEPS_SKIP_REPO_OVERWRITE=true
helm template . -f deps://dummy
```

## Changelog V2

Here are the major changes for version 2:

**Improved Runtime**
The runtime has been reduce significantly by `50-80%` depending on how many chart depedencies need to be installed.
This was achieved by using go routines to install the dependencies in parallel, that are on the same hierarchy level.

This benchmark was done on a chart with 3 hierarchy levels and around 40 dependencies in total.

```bash
helm template . -f deps://dummy -f env/qa/values.yaml  181.52s user 12.70s system 152% cpu 2:07.40 total # V1

export HELM_DEPS_SKIP_REPO_OVERWRITE=true
helm template . -f deps://dummy -f env/qa/values.yaml  200.83s user 19.00s system 358% cpu 1:01.37 total # V2 Default

export HELM_DEPS_SKIP_REFRESH=true
export HELM_DEPS_SKIP_REPO_OVERWRITE=true
export HELM_DEPS_SKIP_START_LOGIN=true
helm template . -f deps://dummy -f env/qa/values.yaml  129.04s user 8.81s system 286% cpu 48.064 total # V2 max performance
```

