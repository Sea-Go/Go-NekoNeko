生成user

```bash
goctl rpc protoc ./proto/user.proto --go_out=./service/user/rpc --go-grpc_out=./service/user/rpc --zrpc_out=./service/user/rpc --style=gozero
```

生成points

```bash
goctl rpc protoc ./proto/points.proto --go_out=./service/points/rpc --go-grpc_out=./service/points/rpc --zrpc_out=./service/points/rpc --style=go_zero
```