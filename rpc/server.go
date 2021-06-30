package rpc

import (
    "fmt"
    "net"
    "google.golang.org/grpc"
)

type (
    Server struct {
        conf RpcConf
        ss   *grpc.Server
        lis  net.Listener
    }
)

func NewServer(rpcConf RpcConf, callbackService func(server *grpc.Server)) *Server {
    lis, err := net.Listen("tcp", fmt.Sprintf("%d", rpcConf.Port))
    if err != nil {
        panic(err)
    }
    ss := grpc.NewServer()
    callbackService(ss)
    return &Server{conf: rpcConf, lis: lis}
}

func (s *Server) Start() {
    //proto.RegisterMsgServer(s, &services.MsgService{})
    err := s.ss.Serve(s.lis)
    if err != nil {
        panic(err)
    }
}

func (s *Server) Stop() {

}
