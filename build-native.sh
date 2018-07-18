#!/usr/bin/env bash

echo "building worker..."
# add `-gcflags '-N -l'` to 'go build ...' to compile for debugging
env GOOS=linux GOARCH=amd64 go build -o workers/golang/golang-worker
echo "worker built"

echo "building samples..."
samples=(
    "HttpTrigger"
    "HttpTriggerBlobBindings"
    "HttpTriggerQueueBindings"
    "HttpTriggerTableBindings"
    "TimerTrigger"
    "BlobTrigger"
    "QueueTrigger"
    "EventGridTrigger"
)

for i in "${samples[@]}"; do
   echo "building $i"
   env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -o "sample/$i/bin/$i.so" "sample/$i/main.go"
   echo "$i built"
done

chmod +rx workers/golang/golang-worker
