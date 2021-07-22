package rpca

import (
    "fmt"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    "github.com/lsamu/ago/rpca/use"
    "google.golang.org/grpc"
    "net"
)

type (
    //Server Server
    Server struct {
        conf ServerConf
        ss   *grpc.Server
        lis  net.Listener
    }
)

//NewServer NewServer
func NewServer(rpcConf ServerConf, callbackService func(server *grpc.Server)) *Server {
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

//AddService 注册服务。。。
func (s *Server) AddService(reg func(s *grpc.Server))  {
    reg(s.ss)
}