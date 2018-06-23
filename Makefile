plugins = grpc
protoc_location = rpc/protos/azure-functions
proto_out_dir = rpc/

GOLANG_WORKER_BINARY = golang-worker
SUBDIRS := $(wildcard sample/*)

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --go_out=plugins=$(plugins):$(proto_out_dir) $(protoc_location)/*.proto

.PHONY: golang-worker
golang-worker:
	GOOS=linux go build -o $(GOLANG_WORKER_BINARY)

.PHONY: dep
dep:
	go get -u github.com/golang/dep/... && \
	dep ensure

.PHONY : samples $(SUBDIRS)
samples : $(SUBDIRS)

$(SUBDIRS) :
	cd $@ && \
	go build -buildmode=plugin
