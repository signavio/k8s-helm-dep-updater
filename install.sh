#!/bin/bash
OWNER=signavio
REPO=k8s-helm-dep-updater
BINARY_NAME=$REPO
FILE_NAME=${BINARY_NAME}_$(uname -s)_$(uname -m).tar.gz

# Get the latest release
RELEASE_RESPONSE=$(curl -sH "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "https://api.github.com/repos/$OWNER/$REPO/releases/latest")

# Extract the id of the asset based on the provided name
ASSET_ID=$(echo "$RELEASE_RESPONSE" | jq --arg name "$FILE_NAME" -r '.assets[] | select(.name==$name) | .id')

# Download the asset file
curl -L -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/octet-stream" \
  "https://api.github.com/repos/$OWNER/$REPO/releases/assets/$ASSET_ID" \
  -o /tmp/$FILE_NAME

tar -xzf /tmp/$FILE_NAME -C /tmp
sudo mv /tmp/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

# HELM_DEP_UPDATER_REPO=k8s-helm-dep-updater
# HELM_DEP_UPDATER_BINARY=$HELM_DEP_UPDATER_REPO
# HELM_DEP_UPDATER_VERSION=latest
# HELM_DEP_UPDATER_FILE_NAME=${HELM_DEP_UPDATER_BINARY}_$(uname -s)_$(uname -m).tar.gz
# HELM_DEP_UPDATER_RELEASE_RESPONSE=$(curl -sH "Authorization: token $GITHUB_TOKEN" -H "Accept: application/vnd.github.v3+json" "https://api.github.com/repos/signavio/$HELM_DEP_UPDATER_REPO/releases/$HELM_DEP_UPDATER_VERSION")
# HELM_DEP_UPDATER_ASSET_ID=$(echo "$HELM_DEP_UPDATER_RELEASE_RESPONSE" | jq --arg name "$HELM_DEP_UPDATER_FILE_NAME" -r '.assets[] | select(.name==$name) | .id')
# curl -L -H "Authorization: token $GITHUB_TOKEN" -H "Accept: application/octet-stream" "https://api.github.com/repos/signavio/$HELM_DEP_UPDATER_REPO/releases/assets/$HELM_DEP_UPDATER_ASSET_ID" -o "/tmp/$HELM_DEP_UPDATER_FILE_NAME"
# tar -xzf "/tmp/$HELM_DEP_UPDATER_FILE_NAME" -C .