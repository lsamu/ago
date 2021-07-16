package rpc

import (
    "fmt"
    "net"
    "google.golang.org/grpc"
)

type (
    //Server Server
    Server struct {
        conf RPCConf
        ss   *grpc.Server
        lis  net.Listener
    }
)
//NewServer NewServer
func NewServer(rpcConf RPCConf, callbackService func(server *grpc.Server)) *Server {
    lis, err := net.Listen("tcp", fmt.Sprintf("%d", rpcConf.Port))
    if err != nil {
        panic(err)
    }
    ss := grpc.NewServer()
    callbackService(ss)
    return &Server{conf: rpcConf, lis: lis}
}

//Start Start
func (s *Server) Start() {
    //proto.RegisterMsgServer(s, &services.MsgService{})
    err := s.ss.Serve(s.lis)
    if err != nil {
        panic(err)
    }
}

//Stop Stop
func (s *Server) Stop() {

}
