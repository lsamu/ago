package sock_io

import (
    "fmt"
    "github.com/gin-gonic/gin"
    socketio "github.com/googollee/go-socket.io"
    "log"
    "net/http"
)

type Server struct {
    conf   SockConf
    engine *gin.Engine
    route  gin.IRoutes
    server *socketio.Server
}

//Route Route
type Route struct {
    Method  string
    Path    string
    Handler gin.HandlerFunc
}

func NewServer(conf SockConf) *Server {
    router := gin.New()
    return &Server{
        engine: router,
        server: socketio.NewServer(nil),
    }
}

func (s *Server) Start() {
    go func() {
        if err := s.server.Serve(); err != nil {
            log.Fatalf("socketio listen error: %s\n", err)
        }
    }()
    defer s.server.Close()
    s.engine.GET("/socket.io/*any", gin.WrapH(s.server))
    s.engine.POST("/socket.io/*any", gin.WrapH(s.server))
    s.engine.StaticFS("/public", http.Dir("../asset"))

    bind := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
    if err := s.engine.Run(bind); err != nil {
        log.Fatal("failed run app: ", err)
    }
}

func (s *Server) Stop() {

}

// Use 中间件
func (e *Server) Use(next gin.HandlerFunc) {
    e.route = e.route.Use(next)
}

// AddRoute 添加路由
func (e *Server) AddRoute(route Route) {
    e.route = e.route.Handle(route.Method, route.Path, route.Handler)
}

// AddEvent 添加事件
func (e *Server) AddEvent(route Route) {
    e.route = e.route.Handle(route.Method, route.Path, route.Handler)
}
