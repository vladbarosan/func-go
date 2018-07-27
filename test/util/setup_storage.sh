#!/usr/bin/env bash

# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# end prolog

account_name=$1
group_name=$2
location=$3

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
echo "storage name $account_name available? $name_available"

echo creating storage account $account_name
account_id=$(az storage account create \
    --name $account_name \
    --resource-group $group_name \
    --location $location \
    --sku 'Standard_LRS' \
    --query id --output tsv)
echo "created storage account: $account_id"

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
