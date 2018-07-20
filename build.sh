#!/usr/bin/env bash

# mount pwd into container and build there
# add `-gcflags '-N -l'` to 'go build ...' to compile for debugging

echo "building worker..."
docker run -it \
    -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker \
    -w /go/src/github.com/Azure/azure-functions-go-worker \
    golang:1.10 /bin/bash -c "go build -o workers/golang/golang-worker"
echo "worker built"

echo "building samples..."
# build samples
samples=(
    "HttpTrigger"
    "HttpTriggerHttpResponse"
    "HttpTriggerBlobBindings"
    "HttpTriggerQueueBindings"
    "HttpTriggerTableBindings"
    "TimerTrigger"
    "BlobTrigger"
    "QueueTrigger"
    "EventGridTrigger"
    "EventHubTriggerEventHubOutput"
    "EventHubTriggerBatchOutput"
)

for i in "${samples[@]}"; do
    echo "building $i"
    docker run -it \
        -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker \
        -w /go/src/github.com/Azure/azure-functions-go-worker \
        golang:1.10 /bin/bash \
            -c "go build -buildmode=plugin -o sample/$i/bin/$i.so sample/$i/main.go"
done

sudo chmod +rx $(pwd)/workers/golang/golang-worker
