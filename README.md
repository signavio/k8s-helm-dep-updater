# Helm Dependency Updater

This project resolves the problem of missing recursive dependency updates in Helm charts. It is a simple tool that can be used to update the dependencies of a Helm chart recursively.

There is an [open issue](https://github.com/helm/helm/issues/2247) in the Helm project to add this functionality to Helm itself. Until this is resolved, this tool can be used to update the dependencies of a Helm chart recursively. It will be part of the [Milestone 3.13.0](https://github.com/helm/helm/milestone/131) which is released on 13th of September 2023.

## Install

```bash
wget -qO- "https://github.com/signavio/k8s-helm-dep-updater/releases/latest/download/k8s-helm-dep-updater_$(uname -s)_$(uname -m).tar.gz" | tar -C /tmp -xzf- k8s-helm-dep-updater
sudo mv /tmp/k8s-helm-dep-updater /usr/local/bin/k8s-helm-dep-updater
```


