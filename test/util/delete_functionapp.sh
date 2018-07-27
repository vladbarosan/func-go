#!/usr/bin/env bash

# prolog
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# end prolog

# config
functionapp_name=$1
functionapp_group_name=$2

functionapp_plan_name="${functionapp_name}-plan"
location=${AZURE_LOCATION_DEFAULT}
# end config

az functionapp delete \
    --name $functionapp_name \
    --resource-group $functionapp_group_name

az appservice plan delete --yes \
    --name $functionapp_plan_name \
    --resource-group $functionapp_group_name
