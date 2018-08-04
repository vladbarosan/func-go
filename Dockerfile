ARG RUNTIME_IMAGE=microsoft/azure-functions-base:dev-nightly
# steps are: 1) build worker, 2) copy worker artifact to runtime image

FROM golang:1.10 as golang-env
WORKDIR /go/src/github.com/vladbarosan/func-go
ENV DEP_RELEASE_TAG=v0.5.0
COPY . .
RUN curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    && dep ensure -v -vendor-only \
    && chmod +x ./build.sh \
    && ./build.sh native

# 3. copy built worker and extensions to runtime image
# ARG instructions used here must be declared before first FROM
FROM ${RUNTIME_IMAGE}

# copy worker to predefined path
COPY --from=golang-env \
    /go/src/github.com/vladbarosan/func-go/workers/golang \
    /azure-functions-host/workers/golang/

# use predefined env var names to point to worker start script
ENV workers:golang:path /azure-functions-host/workers/golang/start.sh
