param(
    $SampleName = $(throw "The sample name is not specified")
)

docker run -it -v ${PWD}:/go/src/github.com/Azure/azure-functions-go-worker -w /go/src/github.com/Azure/azure-functions-go-worker golang:1.10 /bin/bash -c "go build -buildmode=plugin -o sample/${SampleName}/${SampleName}.so sample/${SampleName}/main.go"
