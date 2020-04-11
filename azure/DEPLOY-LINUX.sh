#!/usr/bin/env bash
set -euo pipefail

# variables
RESOURCE_GROUP='200300-hello-gopher'
LOCATION='eastus'
SUBSCRIPTION_ID=$(az account show | jq -r .id)
SCOPE="/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RESOURCE_GROUP}"
# RANDOM_STR='2b7222'
RANDOM_STR=$(echo -n "$SCOPE" | shasum | head -c 6)
STORAGE_NAME="storage${RANDOM_STR}"
FUNCTION_NAME="functions${RANDOM_STR}"
CREATE_IF_EXISTS="false"

# set by actions workflow
# GITHUB_SHA=''
[[ -z "$GITHUB_SHA" ]] && GITHUB_SHA='test'

echo "RANDOM_STR: ${RANDOM_STR}"
echo "GITHUB_SHA: ${GITHUB_SHA}"
return # test

TMP=$(az storage account list -g $RESOURCE_GROUP | jq '[.[].name | index("'$STORAGE_NAME'")] | length')

if [[ "$TMP" == "0" || $CREATE_IF_EXISTS == "true" ]]; then
    echo "az storage account create..."
    az storage account create -g $RESOURCE_GROUP -l $LOCATION -n $STORAGE_NAME \
        --kind StorageV2 \
        --sku Standard_LRS \
        > /dev/null
else
    echo "storage exists..."
fi

TMP=$(az functionapp list -g $RESOURCE_GROUP | jq '[.[].name | index("'$FUNCTION_NAME'")] | length')

if [[ "$TMP" == "0" || $CREATE_IF_EXISTS == "true" ]]; then
    echo "az functionapp create..."
    az functionapp create -g $RESOURCE_GROUP -s $STORAGE_NAME -n $FUNCTION_NAME \
        --consumption-plan-location $LOCATION \
        --os-type Linux \
        --runtime python \
        --functions-version 3 \
        > /dev/null

    echo "az functionapp appsettings..."
    az functionapp config appsettings set -g $RESOURCE_GROUP -n $FUNCTION_NAME \
        --settings "FUNCTIONS_EXTENSION_VERSION=~3" \
        > /dev/null

    az functionapp config appsettings set -g $RESOURCE_GROUP -n $FUNCTION_NAME --settings \
        "WEBSITE_MOUNT_ENABLED=1" \
        "SERVER_NAME=hello-gopher-${GITHUB_SHA}" \
        > /dev/null
else
    echo "functionapp exists..."
fi

echo "build binary..."
source build-container-linux.sh

echo "deploy function..."
cp host.linux.json host.json
source deploy-storage.sh

echo "curl https://${FUNCTION_NAME}.azurewebsites.net/api/healthz"
curl -s -w '%{time_starttransfer}\n' "https://${FUNCTION_NAME}.azurewebsites.net/api/healthz"
