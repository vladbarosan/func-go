#! /bin/bash

# Add -gcflags '-N -l'to 'go build ...' to compile for debugging
docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -o workers/golang/golang-worker"
#Uncomment next lines to build the sample also
docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -buildmode=plugin -o sample/HttpTriggerGo/bin/HttpTriggerGo.so sample/HttpTriggerGo/main.go"
docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -buildmode=plugin -o sample/HttpTriggerBlobBindingsGo/bin/HttpTriggerBlobBindingsGo.so sample/HttpTriggerBlobBindingsGo/main.go"

sudo chmod +rx $(pwd)/workers/golang/golang-worker
