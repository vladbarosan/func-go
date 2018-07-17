#!/usr/bin/env bash
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
source "${__dirname}/.env"
group_name=${AZURE_STORAGE_ACCOUNT_GROUP_NAME}
account_name=${AZURE_STORAGE_ACCOUNT_NAME}
location=${AZURE_LOCATION_DEFAULT}

container_names=(demo demo-out)
queue_name=demoqueue

echo "getting key for account ${account_name}"
key=$(az storage account keys list \
    --account-name $account_name \
    --resource-group $group_name \
    --query '[0].value' -o tsv)

echo -n "putting message in queue, ID: "
az storage message put \
    --content 'hello world' \
    --queue-name $queue_name \
    --account-key $key \
    --account-name $account_name \
    --output tsv --query id

echo -n "getting message from queue, ID: "
az storage message get \
    --queue-name $queue_name \
    --account-key $key \
    --account-name $account_name \
    --output tsv --query id
