#! /bin/bash

# Directory name for start.sh
DIR="$(dirname $0)"
echo "starting the golang worker from $DIR"
/dlv --listen=:40000 --headless=true --api-version=2 exec $DIR/golang-worker -- $@