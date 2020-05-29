IMAGE=token
TAG=$(shell git describe --always)
REGISTRY=registry.cn-shenzhen.aliyuncs.com/pangxieke

default: fmt

fmt:
	go fmt ./...

all:
	docker build -t ${REGISTRY}/${IMAGE}:${TAG} .

push: all
	docker push ${REGISTRY}/${IMAGE}:${TAG}

publish: push
	docker tag ${REGISTRY}/${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:latest
	docker push ${REGISTRY}/${IMAGE}:latest

builder:
	docker build -t ${REGISTRY}/${IMAGE}-builder:latest . -f builder.dockerfile
	docker push ${REGISTRY}/${IMAGE}-builder:latest

doc: spec/*.proto
	docker run --rm -v $(shell pwd)/doc:/out -v $(shell pwd)/spec:/protos pseudomuto/protoc-gen-doc

test:
	gotest -v ./... || go test -v ./...

lint:
	ls -l | grep '^d' | awk '{print $$NF}' | grep -v vender | xargs golint

count:
	cloc --progress=1 ./ --exclude-dir=vendor,doc,pb

pb: pb/token.pb.go

pb/%.pb.go: %.proto
	mkdir -p pb && protoc -I. $^ --go_out=plugins=grpc:pb

.PHONY: test
