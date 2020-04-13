# AZURE

## az group deployment - empty
```bash
az group deployment create --resource-group $RESOURCE_GROUP \
    --template-uri https://raw.githubusercontent.com/asw101/cloud-snips/master/arm/empty/empty.json \
    --mode 'Complete'
```
