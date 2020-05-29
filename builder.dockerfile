FROM golang:latest as builder
RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y git mercurial ca-certificates unzip
# protoc-compiler from apt-get lacks include files, so get a newer one from github
RUN mkdir /protoc \
    && cd /protoc \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip \
    && unzip protoc-3.7.1-linux-x86_64.zip \
    && rm protoc-3.7.1-linux-x86_64.zip
RUN go get github.com/golang/dep/cmd/dep
RUN go get github.com/golang/protobuf/protoc-gen-go
