
ARG NAMESPACE=microsoft
ARG HOST_TAG=dev-nightly

#Install any extension used
FROM microsoft/dotnet:2.1-sdk AS dotnet-env
COPY sample /sample
RUN dotnet build /sample -o bin

FROM golang:1.10 as golang-env

WORKDIR /go/src/github.com/Azure/azure-functions-go-worker
ENV DEP_RELEASE_TAG=v0.4.1
COPY . .
COPY --from=dotnet-env /sample ./sample
RUN  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -vendor-only
RUN chmod +x ./build-native.sh && ./build-native.sh


# Build runtime + worker image
FROM ${NAMESPACE}/azure-functions-base:${HOST_TAG}

# copy the worker in the pre-defined path
COPY --from=golang-env /go/src/github.com/Azure/azure-functions-go-worker/workers/golang /azure-functions-host/workers/golang/
# copy the samples in the pre-defined path
COPY --from=golang-env /go/src/github.com/Azure/azure-functions-go-worker/sample /home/site/wwwroot
ENV workers:golang:path /azure-functions-host/workers/golang/start.sh