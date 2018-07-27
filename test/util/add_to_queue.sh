#!/usr/bin/env bash
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

account_name=$1
group_name=$2

container_names=(demo demo-out)
queue_name=demoqueue

echo "getting key for account ${account_name}"
key=$(az storage account keys list \
    --account-name $account_name \
    --resource-group $group_name \
    --query '[0].value' --output tsv)

echo -n "putting message in queue, ID: "
az storage message put \
    --content 'hello world' \
    --queue-name $queue_name \
    --account-key $key \
    --account-name $account_name \
    --query 'id' --output tsv

echo -n "getting message from queue, ID: "
az storage message get \
    --queue-name $queue_name \
    --account-key $key \
    --account-name $account_name \
    --query '[0].id' --output tsv
echo
