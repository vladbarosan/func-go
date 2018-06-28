#! /bin/bash

docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -gcflags '-N -l' -o workers/golang/golang-worker"
#Uncomment next line to build the sample also
#docker run -it -v $(pwd):/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -buildmode=plugin -gcflags "-N -l" -o sample/HttpTriggerGo/bin/HttpTriggerGo.so sample/HttpTriggerGo/main.go"

sudo chmod +rx $(pwd)/workers/golang/golang-worker
