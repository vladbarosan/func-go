#! /bin/bash

# Add -gcflags '-N -l'to 'go build ...' to compile for debugging
docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -o workers/golang/golang-worker"
#Uncomment next lines to build the sample also
samples=(
    "HttpTrigger"
    "HttpTriggerBlobBindings"
    "HttpTriggerQueueBindings"
    "HttpTriggerTableBindings"
    "TimerTrigger"
    "EventGridTrigger"
    "BlobTrigger"
    "QueueTrigger"
)

for i in "${samples[@]}"
do
   echo "building $i"
   docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -buildmode=plugin -o sample/$i/bin/$i.so sample/$i/main.go"
done

sudo chmod +rx $(pwd)/workers/golang/golang-worker
