# AZURE-FUNCTIONS

See: [AZURE-FUNCTIONS.md](AZURE-FUNCTIONS.md)

## functions (linux)
```bash
az group create -l $LOCATION -n $RESOURCE_GROUP

az storage account create -g $RESOURCE_GROUP -l $LOCATION -n $STORAGE_NAME \
    --kind StorageV2 \
    --sku Standard_LRS

az functionapp create -g $RESOURCE_GROUP -s $STORAGE_NAME -n $FUNCTION_NAME \
    --consumption-plan-location $LOCATION \
    --os-type Linux \
    --runtime python \
    --functions-version 3
```
