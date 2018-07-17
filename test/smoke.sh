#!/usr/bin/env bash
declare stop_containers=${STOP_CONTAINERS:-"0"}
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
source "${__dirname}/.env"

group_name=${AZURE_STORAGE_ACCOUNT_GROUP_NAME}
account_name=${AZURE_STORAGE_ACCOUNT_NAME}
location=${AZURE_DEFAULT_LOCATION}

my_container_name=go-functions-tester

connstr=$(az storage account show-connection-string \
    --name $account_name \
    --resource-group $group_name \
    --key primary \
    --protocol https \
    --query connectionString --output tsv)

echo "ensuring storage account and containers are provisioned"
${__dirname}/setup_storage.sh

echo "building worker into container"
docker build -t azure-functions-go-worker .

echo "running runtime with storage connection"
# running detached and will connect at the end
docker run \
    --detach \
    --name $my_container_name \
    --rm -p 81:80 -e AzureWebJobsStorage="$connstr" \
    azure-functions-go-worker

echo
echo "creating and deleting a blob"
${__dirname}/upload_blob.sh

echo
echo "sending and receiving a queue message"
${__dirname}/add_to_queue.sh

if [[ "$stop_containers" != "0" ]]; then
    echo "shutting down container"
    docker stop $my_container_name
else 
    echo "leaving container running and attaching"
    docker attach $my_container_name
fi
