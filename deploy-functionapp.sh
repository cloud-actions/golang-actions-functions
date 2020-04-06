[[ -z "$RESOURCE_GROUP" ]] && echo 'RESOURCE_GROUP not set!' && return
[[ -z "$FUNCTION_NAME" ]] && echo 'FUNCTION_NAME not set!' && return
echo "RESOURCE_GROUP: ${RESOURCE_GROUP}"
echo "FUNCTION_NAME: ${FUNCTION_NAME}"

mkdir -p _/

FILE_NAME='_/deploy.zip'
zip -r $FILE_NAME .

az functionapp deployment source config-zip \
    -g $RESOURCE_GROUP -n $FUNCTION_NAME \
    --src $FILE_NAME

rm $FILE_NAME
