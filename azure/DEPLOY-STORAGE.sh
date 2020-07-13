[[ -z "$RESOURCE_GROUP" ]] && echo 'RESOURCE_GROUP not set!' && return
[[ -z "$FUNCTION_NAME" ]] && echo 'FUNCTION_NAME not set!' && return
echo "RESOURCE_GROUP: ${RESOURCE_GROUP}"
echo "FUNCTION_NAME: ${FUNCTION_NAME}"

mkdir -p _/

FILE_NAME='_/deploy.zip'
zip -r $FILE_NAME .

export AZURE_STORAGE_CONNECTION_STRING=$(az functionapp config appsettings list -g $RESOURCE_GROUP -n $FUNCTION_NAME | jq -r '.[] | select(.name == "AzureWebJobsStorage").value')

CONTAINER_NAME='function-releases'
BLOB_UPLOAD="${FUNCTION_NAME}-primary.zip"
BLOB_DELETE="${FUNCTION_NAME}-secondary.zip"

SAS_EXPIRY=$(($(date +"%Y")+10))'-01-01'

echo "az storage container create..."
az storage container create -n $CONTAINER_NAME

BLOB_UPLOAD_COUNT=$(az storage blob list -c $CONTAINER_NAME | jq '[.[] | select(.name == "'$BLOB_UPLOAD'")] | length')
if [[ "$BLOB_UPLOAD_COUNT" != "0" ]]; then
    TMP=$BLOB_UPLOAD
    BLOB_UPLOAD=$BLOB_DELETE
    BLOB_DELETE=$TMP
fi

echo "az storage blob upload (${BLOB_UPLOAD})..."

az storage blob upload -c $CONTAINER_NAME -n $BLOB_UPLOAD -f $FILE_NAME

BLOB_URL=$(az storage blob url -c $CONTAINER_NAME -n $BLOB_UPLOAD | jq -r '.')
BLOB_SAS=$(az storage blob generate-sas -c $CONTAINER_NAME -n $BLOB_UPLOAD \
    --permissions acdrw --expiry $SAS_EXPIRY | jq -r '.')

echo "az functionapp config appsettings..."
az functionapp config appsettings set -g $RESOURCE_GROUP -n $FUNCTION_NAME \
    --settings "WEBSITE_RUN_FROM_PACKAGE=${BLOB_URL}?${BLOB_SAS}" \
    > /dev/null

BLOB_DELETE_COUNT=$(az storage blob list -c $CONTAINER_NAME | jq '[.[] | select(.name == "'$BLOB_DELETE'")] | length')
if [[ "$BLOB_DELETE_COUNT" != "0" ]]; then
    echo "az storage blob delete (${BLOB_DELETE})..."
    az storage blob delete -c $CONTAINER_NAME -n $BLOB_DELETE
fi

rm $FILE_NAME
