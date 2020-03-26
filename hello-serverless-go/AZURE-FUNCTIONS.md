# AZURE-FUNCTIONS

## variables
```bash
RESOURCE_GROUP='200225-functions'
LOCATION='eastus'
RANDOM_STR='14f0ce'
if [ -z "$RANDOM_STR" ]; then RANDOM_STR=$(openssl rand -hex 3); else echo $RANDOM_STR; fi
STORAGE_NAME="storage${RANDOM_STR}"
FUNCTION_NAME="functions${RANDOM_STR}"
```

## functions
```bash
az group create -l $LOCATION -n $RESOURCE_GROUP

az storage account create -g $RESOURCE_GROUP -l $LOCATION -n $STORAGE_NAME \
    --kind StorageV2 \
    --sku Standard_LRS

az functionapp create -g $RESOURCE_GROUP -s $STORAGE_NAME -n $FUNCTION_NAME \
    --consumption-plan-location $LOCATION \
    --os-type Windows \
    --runtime dotnet \
    --functions-version 3

az functionapp config set -g $RESOURCE_GROUP -n $FUNCTION_NAME \
    --use-32bit-worker-process false
```

## empty
```bash
az group deployment create --resource-group $RESOURCE_GROUP \
    --template-uri https://raw.githubusercontent.com/asw101/cloud-snips/master/arm/empty/empty.json \
    --mode 'Complete'
```

## deploy
```bash
RESOURCE_GROUP='200225-functions'
RANDOM_STR='14f0ce'
FUNCTION_NAME="functions${RANDOM_STR}"
cp host.windows.json host.json
source deploy.sh
cp host.linux.json host.json
```
