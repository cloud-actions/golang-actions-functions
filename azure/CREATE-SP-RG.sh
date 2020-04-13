#!/usr/bin/env bash
set -euo pipefail

[[ -z "${RESOURCE_GROUP:-}" ]] && RESOURCE_GROUP='200400-hello-github'
[[ -z "${LOCATION:-}" ]] && LOCATION='eastus'

SUBSCRIPTION_ID=$(az account show | jq -r .id)
SCOPE="/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RESOURCE_GROUP}"
RANDOM_STR=$(echo -n "$SCOPE" | shasum | head -c 6)

az group create -n $RESOURCE_GROUP -l $LOCATION

SP=$(az ad sp create-for-rbac --sdk-auth -n "${RESOURCE_GROUP}-${RANDOM_STR}" --role contributor \
    --scopes $SCOPE)

echo $SP | jq -c
