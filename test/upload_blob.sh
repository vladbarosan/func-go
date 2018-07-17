#!/usr/bin/env bash
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
source "${__dirname}/.env"
group_name=${AZURE_STORAGE_ACCOUNT_GROUP_NAME}
account_name=${AZURE_STORAGE_ACCOUNT_NAME}
location=${AZURE_LOCATION_DEFAULT}

# shared with setup_storage.sh
container_names=(demo demo-out)
blob_name=blob.bin
blob_path=${__dirname}/testdata/${blob_name}

echo "getting key for account ${account_name}"
key=$(az storage account keys list \
    --account-name $account_name \
    --resource-group $group_name \
    --query '[0].value' -o tsv)

echo "uploading blob $blob_path"
az storage blob upload \
    --container-name ${container_names[0]} \
    --file $blob_path \
    --name $blob_name \
    --account-key $key \
    --account-name $account_name \
    --output json --no-progress > /dev/null

echo "deleting blob ${container_names[0]}/${blob_name}"
az storage blob delete \
    --container-name ${container_names[0]} \
    --name $blob_name \
    --account-key $key \
    --account-name $account_name \
    --output json > /dev/null
