#!/usr/bin/env bash
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

source "${__dirname}/.env"
group_name=${AZURE_STORAGE_ACCOUNT_GROUP_NAME}
account_name=${AZURE_STORAGE_ACCOUNT_NAME}
location=${AZURE_LOCATION_DEFAULT}

container_names=(demo demo-out)
queue_name=demoqueue

echo creating group $group_name
az group create \
    --name $group_name \
    --location $location \
    --output tsv --query id

name_available=$(az storage account check-name \
    -n $account_name \
    --query nameAvailable --output tsv)
# if not name_available: account_name += 01

echo creating storage account $account_name
account_id=$(az storage account create \
    --name $account_name \
    --resource-group $group_name \
    --location $location \
    --query id --output tsv)

echo getting account key
key=$(az storage account keys list \
    --account-name $account_name \
    --resource-group $group_name \
    --query '[0].value' -o tsv)

echo getting account connstr
connstr=$(az storage account show-connection-string \
            --ids $account_id \
            --query connectionString --output tsv)

echo creating containers
for container_name in ${container_names[@]}; do
    az storage container create \
        --name $container_name \
        --account-key $key \
        --account-name $account_name \
        --output tsv --query name
done

echo creating queues
az storage queue create \
    --name $queue_name \
    --account-key $key \
    --account-name $account_name \
    --query name --output tsv
