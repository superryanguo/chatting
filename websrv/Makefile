
GOPATH:=$(shell go env GOPATH)
MODIFY=Mgithub.com/micro/go-micro/api/proto/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
    
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/websrv/websrv.proto
    

.PHONY: build
build: proto

	go build -o websrv-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t websrv-service:latest

.PHONY: mache
mache:
	protoc --go_out=. proto/mache/mache.proto
