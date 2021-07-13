# ago
```
基于go+gin+gorm+grpc+socket开发框架. 持续更新中...
```

## 安装
```
go get -u github.com/lsamu/ago
```

## http
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

```

## socket
```

```

## tcp
```

```