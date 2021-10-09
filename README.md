# ago
```
基于go+gin+gorm+grpc+socket开发框架. 持续更新中...
```

## 安装
```
go get -u github.com/lsamu/ago
```

## rest
```
func main() {
    server:=rest.NewServer(rest.RestConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()
    server.Start()
}
```

## rpc
```
func main() {
    server:= rpc.NewServer(rpc.ServerConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    server.AddService(func(s *grpc.Server) {
        //注册服务
    })
    defer server.Stop()
    server.Start()
}
```

## socketio
```
func main() {
    server:= sockio.NewServer(sockio.SockConf{
        Host: "0.0.0.0",
        Port: 8888,
    })
    defer server.Stop()
    server.Start()
}

```

## tcp
```

```