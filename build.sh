#!/usr/bin/env bash

# mount pwd into container and build there
# add `-gcflags '-N -l'` to 'go build ...' to compile for debugging
docker run -it \
    -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker \
    -w /go/src/github.com/Azure/azure-functions-go-worker \
    golang:1.10 /bin/bash -c "go build -o workers/golang/golang-worker"

# build samples
samples=(
    "HttpTrigger"
    "HttpTriggerHttpResponse"
    "HttpTriggerBlobBindings"
    "HttpTriggerQueueBindings"
    "HttpTriggerTableBindings"
    "TimerTrigger"
    "EventGridTrigger"
    "BlobTrigger"
    "QueueTrigger"
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
