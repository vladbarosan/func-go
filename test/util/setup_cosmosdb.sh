#!/usr/bin/env bash

# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# end prolog

account_name=$1
group_name=$2
location=$3

db_name="Documents"
collection_names=(reports tasks)

"echo creating group $group_name"
az group create \
    --name $group_name \
    --location $location \
    --output tsv --query id

name_available=$(az cosmosdb check-name-exists \
    --name $account_name \
     --output tsv)
echo "cosmos db account name $account_name available? $name_available"

echo "creating cosmosdb account $account_name"
account_id=$(az cosmosdb create \
    --name $account_name \
    --resource-group $group_name \
    --kind 'GlobalDocumentDB' \
    --query id --output tsv)
echo "created cosmosdb account: $account_id"


echo "getting account key"
account_key=$(az cosmosdb list-keys \
    --name $account_name \
    --resource-group $group_name \
    --query "primaryMasterKey" -o tsv)

connstr="AccountEndpoint=https://$account_name.documents.azure.com:443/;AccountKey=$account_key;"

echo "creating cosmosdb database $db_name"
az cosmosdb database create \
    --db-name $db_name \
    --name $account_name \
    --resource-group $group_name

echo "creating collections"
for collection_name in ${collection_names[@]}; do
    az cosmosdb collection create \
        --collection-name $collection_name \
        --db-name $db_name
        --resource-group $group_name
        --output tsv --query name
done
