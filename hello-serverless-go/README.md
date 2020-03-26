# README

## hello-serverless-go/
```bash
cd hello-serverless-go/

# build & copy
cd ../hello-gopher/
source build.sh
cd ../hello-serverless-go/

# run
docker run -p 7071:7071 -v ${PWD}:/pwd/ -w /pwd/ --rm -it func:latest bash
func start

# test
curl http://localhost:7071/api/HttpTrigger?name=Aaron

echo '{"name": "Aaron"}' | curl -d @- http://localhost:7071/api/HttpTrigger?name=Aaron

curl http://localhost:7071/api/HttpTrigger2?name=Aaron

echo '{"name": "Aaron"}' | curl -d @- http://localhost:7071/api/HttpTrigger2
```

## storage account for timer trigger
```bash
STORAGE_ACCOUNT='storageaccount2002090c6'
STORAGE_CONNECTION_STRING=$(az storage account show-connection-string --name $STORAGE_ACCOUNT | jq -r '.connectionString')
# echo "func settings add AzureWebJobsStorage '${STORAGE_CONNECTION_STRING}'"
func settings add AzureWebJobsStorage $STORAGE_CONNECTION_STRING
```

## deploy

See: [AZURE-FUNCTIONS.md](AZURE-FUNCTIONS.md#deploy)

## test
```bash
FUNCTION_NAME='functions9522be'

# curl
curl "https://${FUNCTION_NAME}.azurewebsites.net/api/HttpTrigger"

curl "https://${FUNCTION_NAME}.azurewebsites.net/api/HttpTrigger?name=Aaron"

# hey
# https://github.com/rakyll/hey
# brew install hey

hey "https://${FUNCTION_NAME}.azurewebsites.net/api/HttpTrigger?name=Aaron"

hey -c 50 -n 2000 "https://${FUNCTION_NAME}.azurewebsites.net/api/HttpTrigger?name=Aaron"

hey -c 100 -n 2000 "https://${FUNCTION_NAME}.azurewebsites.net/api/HttpTrigger?name=Aaron"
```
