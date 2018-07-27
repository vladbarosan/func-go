#!/usr/bin/env bash
# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# end prolog

# global_config
# TODO(joshgav): replace with single ID
declare account_name=${1}
declare group_name=${2}
# end global_config

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
