#!/usr/bin/env bash

# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
source "${__dirname}/../../.env"
# end prolog

# config
functionapp_name=$1
functionapp_group_name=$2
sa_name=$3
sa_group_name=$4

functionapp_plan_name="${functionapp_name}-plan"
location=${AZURE_LOCATION_DEFAULT}
runtime_image_uri="${RUNTIME_IMAGE_REGISTRY}/${RUNTIME_IMAGE_REPO}:${RUNTIME_IMAGE_TAG}"
# end config

function ensure_group () {
    local group_name=$1
    local group_id

    group_id=$(az group show --name $group_name --query 'id' --output tsv 2> /dev/null)

    if [[ -z $group_id ]]; then
        group_id=$(az group create \
            --name $group_name \
            --location $location \
            --query 'id' --output tsv)
    fi
    echo $group_id
}

ensure_group $functionapp_group_name

echo "get storage account"
storage_account_id=$(az storage account show \
    --name $sa_name \
    --resource-group $sa_group_name \
    --query 'id' --output tsv)
echo "storage account: $storage_account_id"

echo "create functionapp plan"
appservice_plan_id=$(az appservice plan create \
    --name $functionapp_plan_name \
    --resource-group $functionapp_group_name \
    --location $location \
    --is-linux \
    --query 'id' --output tsv)
echo "functionapp plan: $appservice_plan_id"

echo "create functionapp"
# CLI command adds lots of config for us :)
# review for troubleshooting: https://github.com/Azure/azure-cli/blob/dev/src/command_modules/azure-cli-appservice/azure/cli/command_modules/appservice/custom.py#create_function
# fixed link: https://github.com/Azure/azure-cli/blob/1558b74ba787b481c1791abd68ed5608d87cee02/src/command_modules/azure-cli-appservice/azure/cli/command_modules/appservice/custom.py#L1657
az functionapp create \
    --name $functionapp_name \
    --resource-group $functionapp_group_name \
    --storage-account $storage_account_id \
    --plan $appservice_plan_id \
    --deployment-container-image-name "${runtime_image_uri}" \
    --output json

