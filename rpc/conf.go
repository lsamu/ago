package rpc

import (
    "fmt"
    "google.golang.org/grpc"
    "log"
    "net"
)

type (
    RpcConf struct {
        Host string
        Port int
    }
)

/**
 * @Description: 启动grpc服务 并注册服务
 * @param port 端口
 * @param callbackService 回调注册服务
 */
func StartServer(port string, callbackService func(server *grpc.Server)) {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        panic(err)
    }
    s := grpc.NewServer()
    callbackService(s)
    //proto.RegisterMsgServer(s, &services.MsgService{})
    err = s.Serve(lis)
    if err != nil {
        panic(err)
    }
}

/**
 * @Description: 创建grpc客户端
 * @param ip
 * @param port
 * @return conn 链接
 * @return err
 */
func StartClient(ip, port string) (conn *grpc.ClientConn, err error) {
    target := fmt.Sprintf("%s:%s", ip, port)
    conn, err = grpc.Dial(target, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
        return nil, err
    }
    defer conn.Close()
    //client := proto.NewMsgClient(conn)
    return conn, err
}
