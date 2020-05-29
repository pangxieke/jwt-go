# 令牌服务
使用jwt搭建令牌服务

## Test
先生成pb文件
`make pb`

启动测试
`make test`

## 启动服务
```
export SECRET=dc3dc8a96e7053c54ee5267363f9cd803912ea82
go run cmd/main.go
```

## 启动客户端
```
go run client/main.go
```