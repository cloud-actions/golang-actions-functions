#!/usr/bin/env bash
set -euo pipefail

# variables
RESOURCE_GROUP='200200-hello-gopher'
LOCATION='eastus'
SUBSCRIPTION_ID=$(az account show | jq -r .id)
SCOPE="/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RESOURCE_GROUP}"
# RANDOM_STR='2b7222'
RANDOM_STR=$(echo -n "$SCOPE" | shasum | head -c 6)
STORAGE_NAME="storage${RANDOM_STR}"
FUNCTION_NAME="functions${RANDOM_STR}"
CREATE_IF_EXISTS="false"

echo "RANDOM_STR: {$RANDOM_STR}"

TMP=$(az storage account list -g $RESOURCE_GROUP | jq '.[].name | index("'$STORAGE_NAME'")')

if [[ "null" == "$TMP" || $CREATE_IF_EXISTS == "true" ]]; then
    echo "az storage account create..."
    az storage account create -g $RESOURCE_GROUP -l $LOCATION -n $STORAGE_NAME \
        --kind StorageV2 \
        --sku Standard_LRS \
        > /dev/null
else
    echo "storage exists..."
fi

TMP=$(az functionapp list -g $RESOURCE_GROUP | jq '.[].name | index("'$FUNCTION_NAME'")')

if [[ "null" == "$TMP" || $CREATE_IF_EXISTS == "true" ]]; then
    echo "az functionapp create..."
    az functionapp create -g $RESOURCE_GROUP -s $STORAGE_NAME -n $FUNCTION_NAME \
        --consumption-plan-location $LOCATION \
        --os-type Windows \
        --runtime dotnet \
        > /dev/null

    echo "az functionapp appsettings..."
    az functionapp config appsettings set -g $RESOURCE_GROUP -n $FUNCTION_NAME \
        --settings "FUNCTIONS_EXTENSION_VERSION=~3" \
        > /dev/null

    echo "az functionapp config..."
    az functionapp config set -g $RESOURCE_GROUP -n $FUNCTION_NAME \
        --use-32bit-worker-process false \
        > /dev/null
else
    echo "functionapp exists..."
fi

echo "build binary..."
cd hello-gopher/
source build-container.sh
cd ..

echo "deploy function..."
cd hello-serverless-go/
cp host.windows.json host.json
source deploy.sh
cd ..

# output URL, etc...
echo "Functions deployed to: https://${FUNCTION_NAME}.azurewebsites.net/"
