[[ -z "$RESOURCE_GROUP" ]] && echo 'RESOURCE_GROUP not set!' && return
[[ -z "$FUNCTION_NAME" ]] && echo 'FUNCTION_NAME not set!' && return
echo "RESOURCE_GROUP: ${RESOURCE_GROUP}"
echo "FUNCTION_NAME: ${FUNCTION_NAME}"

mkdir -p _/
zip -r _/deploy.zip .

az functionapp deployment source config-zip \
    -g $RESOURCE_GROUP -n $FUNCTION_NAME \
    --src _/deploy.zip

rm _/deploy.zip
