
ARG NAMESPACE=microsoft
ARG HOST_TAG=2.0

#start from golang 1.10 (as multiple plugins with same package name fails in golang-1.9.x)
FROM golang:1.10 as golang-env

WORKDIR /go/src/github.com/Azure/azure-functions-go-worker
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -vendor-only

RUN go build -o golang/golang-worker

FROM ${NAMESPACE}/azure-functions-base:${HOST_TAG}

# copy the worker in the pre-defined path
COPY --from=golang-env /go/src/github.com/Azure/azure-functions-go-worker/golang /azure-functions-host/workers/golang/