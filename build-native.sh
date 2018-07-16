#! /bin/bash

echo "building worker..."

# Add -gcflags '-N -l'to 'go build ...' to compile for debugging
env GOOS=linux GOARCH=amd64 go build -o workers/golang/golang-worker
echo "worker built"
echo "building samples..."

#Uncomment next lines to build the sample also
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

for i in "${samples[@]}"
do
   echo "building $i"
   env GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o "sample/$i/bin/$i.so" "sample/$i/main.go"
   echo "$i built"
done

chmod +rx workers/golang/golang-worker