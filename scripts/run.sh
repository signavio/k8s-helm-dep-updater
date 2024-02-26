#!/bin/sh

# Make sure k8s-helm-dep-updater is in PATH
command -v k8s-helm-dep-updater >/dev/null 2>&1 || { echo >&2 "k8s-helm-dep-updater not found. Please install it first."; exit 1; }

k8s-helm-dep-updater >/dev/null