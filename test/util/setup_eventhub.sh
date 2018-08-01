#!/usr/bin/env bash

# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# end prolog

namespace_name=$1
group_name=$2
location=$3

eventhub_names=(demo-go-func-in demo-go-func-out demo-go-func-batch-in)

echo creating group $group_name
az group create \
    --name $group_name \
    --location $location \
    --output tsv --query id

name_available=$(az eventhubs namespace exists \
    -n $namespace_name \
    --query nameAvailable --output tsv)
echo "event hub namespace name $namespace_name available? $name_available"

echo "creating event hub namespace $namespace_name"
namespace_id=$(az eventhubs namespace create \
    --name $namespace_name \
    --resource-group $group_name \
    --location $location \
    --sku 'Standard' \
    --query id --output tsv)
echo "created event hub namespace: $namespace_id"

echo getting namespace default SAS Policy connection string
policy_name=$(az eventhubs namespace authorization-rule list \
    --namespace-name $namespace_name \
    --resource-group $group_name \
    --query "[0].name" -o tsv)
connstr=$(az eventhubs namespace authorization-rule keys list \
    --namespace-name $namespace_name
    --resource-group $group_name
    --name $policy_name
    --query "primaryConnectionString" -o tsv)

echo creating event hubs
for eventhub_name in ${eventhub_names[@]}; do
    az eventhubs eventhub create \
        --name $eventhub_name \
        --namespace-name $namespace_name
        --resource-group $group_name
        --output tsv --query name
done
