
ARG NAMESPACE=microsoft
ARG HOST_TAG=2.0
ARG MODE="dev"

#start from golang 1.10 (as multiple plugins with same package name fails in golang-1.9.x)
FROM golang:1.10 as golang-env

WORKDIR /go/src/github.com/Azure/azure-functions-go-worker
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -vendor-only

# The -gcflags "all=-N -l" flag helps us get a better debug experience
RUN go build -gcflags "-N -l" -o golang/golang-worker

# This is prod
#RUN go build -o golang/golang-worker

# Compile Delve
RUN go get github.com/derekparker/delve/cmd/dlv

FROM ${NAMESPACE}/azure-functions-base:${HOST_TAG}

# copy the worker in the pre-defined path
COPY --from=golang-env /go/src/github.com/Azure/azure-functions-go-worker/golang /azure-functions-host/workers/golang/

COPY --from=golang-env /go/bin/dlv /