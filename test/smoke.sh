#!/usr/bin/env bash

## prolog
set -o errexit
__filename=${BASH_SOURCE[0]}
__dirname=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
__root=$(cd "${__dirname}/../" && pwd)
if [[ ! -f "${__root}/.env" ]]; then cp "${__root}/.env.tpl" "${__root}/.env"; fi
source "${__root}/.env"
## end prolog

## parameters
declare -i publish=${1:-0}  # 0: false; 1: true
declare -i hosted=${2:-0}   # 0: false; 1: true
declare -i continue=${3:-1} # 0: stop;  1: follow; 2: background
declare run_image_uri=${4:-"${RUNTIME_IMAGE_REGISTRY}/${RUNTIME_IMAGE_REPO}:${RUNTIME_IMAGE_TAG}"}
declare instance_name=${5:-${RUNTIME_INSTANCE_NAME}}
declare sa_name="${6:-${STORAGE_ACCOUNT_NAME}smoker}"
declare group_name="${7:-${AZURE_GROUP_NAME_BASE}-smoker}"
declare location=${8:-${AZURE_LOCATION_DEFAULT}}
## end parameters

## prepare run image
echo "building run image [${run_image_uri}]"
${__dirname}/build.sh $publish
## end prepare run image

## ensure storage
echo "ensuring storage account [${sa_name}]"
${__dirname}/util/setup_storage.sh $sa_name $group_name $location
## end ensure storage

## test instance
# start_instance starts an instance of the Functions runtime
# :param: hosted int 0: 0: local Docker; 1: App Service Functions
function start_instance() {
    declare -i hosted=${1:-$hosted}  # 0: local Docker; 1: App Service plan
    declare instance_name=${2:-${instance_name}}
    declare sa_name=${3:-${sa_name}}
    declare sa_group_name=${4:-${group_name}}
    declare published_port=${5:-${RUNTIME_INSTANCE_PORT}} # only used for local

    if [[ $hosted == 1 ]]; then
        echo "running runtime in App Service plan"
        # workaround for <https://github.com/Azure/azure-cli/issues/6918>
        # TODO: when resolved, don't require a group name at all, set that up
        #       automatically
        instance_group_name=$sa_group_name
        # TODO: check $instance_name availability
        ${__dirname}/util/setup_functionapp.sh \
            $instance_name $instance_group_name $sa_name $sa_group_name
    else
        echo "running runtime locally via Docker"
        # stop current running instance if necessary
        cid=$(docker container ls --filter "name=$instance_name" --quiet)
        if [[ -n $cid ]]; then
            docker container stop $instance_name
        fi

        connstr=$(az storage account show-connection-string \
            --name $sa_name \
            --resource-group $sa_group_name \
            --key primary \
            --protocol https \
            --query 'connectionString' --output tsv)

        # running detached first so we can run tests
        # will attach or stop at end per $continue
        docker container run --rm --detach \
            --name $instance_name \
            --publish "${published_port}:80" \
            --env "AzureWebJobsStorage=$connstr" \
            "${run_image_uri}"

        # wait for worker to be ready
        # TODO: queue requests in worker till ready
        sleep 10
    fi
}

# continue_instance attaches to or deletes an instance of the Functions runtime
# :param: hosted int 0: 0: local Docker; 1: App Service Functions
function continue_instance () {
    declare -i hosted=${1:-0}  # 0: local Docker; 1: App Service plan
    declare -i continue=${2:-1}
    declare instance_name=${3:-${instance_name}}

    if [[ $hosted == 1 ]]; then
        # workaround for <https://github.com/Azure/azure-cli/issues/6918>
        instance_group_name=$sa_group_name
        if [[ $continue == 1 ]]; then
            echo "leaving hosted instance running and following logs"
            ${__dirname}/util/connect_functionapp.sh $instance_name $instance_group_name
        elif [[ $continue == 2 ]]; then
            echo "leaving hosted instance running"
        else
            echo "deleting hosted instance"
            ${__dirname}/util/delete_functionapp.sh $instance_name $instance_group_name
        fi
    else # [[ $hosted != 1 ]]
        if [[ $continue == 1 ]]; then
            # this is the default set in parameter declaration
            echo "leaving local instance running and following logs"
            echo -e "run \`docker container stop ${instance_name}\` to stop and remove\n"
            docker container logs --follow $instance_name
        elif [[ $continue == 2 ]]; then
            echo "leaving local instance running"
            echo "run \`docker container logs [--follow] ${instance_name}\` to get [and follow] logs"
            echo "run \`docker container stop ${instance_name}\` to stop and remove"
        else
            echo "shutting down and removing local instance"
            docker container stop $instance_name
        fi
    fi
}


# test_instance invokes tests against an instance of the Functions runtime
function test_instance () {
    declare -i hosted=${1:-0}  # 0: local Docker; 1: App Service plan
    declare instance_name=${2:-${instance_name}}
    declare instance_domain=${3:-${RUNTIME_INSTANCE_DOMAIN}}
    declare instance_port=${4:-${RUNTIME_INSTANCE_PORT}}

    if [[ $hosted == 1 ]]; then
        test_hostname=${instance_name}.${instance_domain}
    else # [[ $hosted != 1 ]]
        test_hostname=localhost:${instance_port}
    fi

    # trigger HttpTrigger
    echo "Test: HttpTrigger"
    person_name="world"
    url="http://${test_hostname}/api/HttpTrigger?name=${person_name}"
    body='{"greeting":"Where would you like to go today?"}'
    echo "POST ${url}; Body: \"${body}\""
    echo -n "Response: "
    curl -L "$url" \
        --data "$body" \
        --header 'Content-Type: application/json'
    # TODO: verify
    echo ""
    echo "end Test: HttpTrigger"
    echo ""

    # trigger BlobTrigger
    echo "Test: BlobTrigger"
    echo "creating and deleting a blob to trigger events"
    ${__dirname}/util/upload_blob.sh $sa_name $group_name
    # TODO: verify
    echo "end Test: BlobTrigger"
    echo ""

    # trigger QueueTrigger
    echo "Test: QueueTrigger"
    echo "sending and receiving a queue message to trigger events"
    ${__dirname}/util/add_to_queue.sh $sa_name $group_name
    # TODO: verify
    echo "end Test: QueueTrigger"
    echo ""
}

start_instance $hosted $instance_name $sa_name $group_name
test_instance $hosted $instance_name
continue_instance $hosted $continue $instance_name
## end test instance
