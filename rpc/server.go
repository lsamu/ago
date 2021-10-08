package rpc

import (
    "fmt"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/lsamu/ago/rpc/use"
    "google.golang.org/grpc"
    "net"
)

type (
    //Server Server
    Server struct {
        conf   ServerConf
        server *grpc.Server
        listen net.Listener
    }
)

//NewServer NewServer
func NewServer(rpcConf ServerConf) *Server {
    lis, err := net.Listen("tcp", fmt.Sprintf("%d", rpcConf.Port))
    if err != nil {
        panic(err)
    }

    //注册中间件
    ss := grpc.NewServer(
        //grpc.StreamInterceptor(
        //    grpc_middleware.ChainStreamServer(
        //        use.Auth(),
        //    )),
        grpc.UnaryInterceptor(
            grpc_middleware.ChainUnaryServer(
                use.Auth2(),
            )))

    return &Server{conf: rpcConf, listen: lis, server: ss}
}

//Start Start
func (s *Server) Start() {
    //proto.RegisterMsgServer(s, &services.MsgService{})
    err := s.server.Serve(s.listen)
    if err != nil {
        panic(err)
    }
}

//Stop Stop
func (s *Server) Stop() {
    s.server.Stop()
}

//Use Use
func (s *Server) Use(handler interface{}) {

}

//GetServer GetServer
func (s *Server) GetServer() *grpc.Server {
    return s.server
}

//AddService 注册服务。。。
func (s *Server) AddService(regService func(s *grpc.Server)) {
    regService(s.server)
}
