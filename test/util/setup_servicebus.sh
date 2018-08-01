#!/usr/bin/env bash

# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# end prolog

namespace_name=$1
group_name=$2
location=$3

queue_names=(inputqueue outputqueue)

echo creating group $group_name
az group create \
    --name $group_name \
    --location $location \
    --output tsv --query id

name_available=$(az servicebus namespace exists \
    -n $namespace_name \
    --query nameAvailable --output tsv)
echo "service bus namespace name $namespace_name available? $name_available"

echo "creating service bus namespace $namespace_name"
namespace_id=$(az servicebus namespace create \
    --name $namespace_name \
    --resource-group $group_name \
    --location $location \
    --sku 'Standard' \
    --query id --output tsv)
echo "created service bus namespace: $namespace_id"

echo getting namespace default SAS Policy connection string
policy_name=$(az servicebus namespace authorization-rule list \
    --namespace-name $namespace_name \
    --resource-group $group_name \
    --query "[0].name" -o tsv)
connstr=$(az servicebus namespace authorization-rule keys list \
    --namespace-name $namespace_name
    --resource-group $group_name
    --name $policy_name
    --query "primaryConnectionString" -o tsv)

echo creating queues
for queue_name in ${queue_names[@]}; do
    az servicebus queue create \
        --name $queue_name \
        --namespace-name $namespace_name
        --resource-group $group_name
        --output tsv --query name
done
