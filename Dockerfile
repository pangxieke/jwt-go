FROM xxx/pangxieke/token-builder:latest as builder

WORKDIR /go/src/token

COPY token.proto .
RUN mkdir pb \
    && PATH=/protoc/bin:$PATH protoc -I/protoc/include -I. --go_out=plugins=grpc:pb token.proto

COPY . .
RUN go build -ldflags '-extldflags "-static"' ./cmd/token

FROM scratch
LABEL description="令牌服务：为业务服务的JWT鉴权令牌" \
      maintainer="pangxieke" \
      SERVICE_NAME=token \
      SERVICE_50051_NAME=token_grpc

EXPOSE 8080
COPY --from=builder /go/src/token/token /bin/app
ENTRYPOINT ["app"]
