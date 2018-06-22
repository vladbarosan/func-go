#! /bin/bash

# Directory name for start.sh
DIR="$(dirname $0)"
echo "starting the golang worker from $DIR"
$DIR/golang-worker $@