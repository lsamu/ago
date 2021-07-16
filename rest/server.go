package rest

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
)

import "github.com/gin-gonic/gin"

type (
    //Server Server
    Server struct {
        conf   RestConf
        engine *gin.Engine
        server *http.Server
        route  gin.IRoutes
    }
    //Route Route
    Route struct {
        Method  string
        Path    string
        Handler gin.HandlerFunc
    }
)

// NewServer 服务
func NewServer(conf RestConf) *Server {
    engine := gin.Default()
    return &Server{
        conf:   conf,
        engine: engine,
        route:  engine.Group(""),
    }
}

// Start 启动
func (e *Server) Start() {
    bind := fmt.Sprintf("%s:%d", e.conf.Host, e.conf.Port)
    e.server = &http.Server{
        Addr:           bind,
        Handler:        e.engine,
        ReadTimeout:    30 * time.Second,
        WriteTimeout:   30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    time.Sleep(10 * time.Microsecond)
    go func() {
        log.Printf("Listening and serving HTTP on %s\n", bind)
        if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutdown Server ...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := e.server.Shutdown(ctx); err != nil {
        log.Fatal("Server Shutdown:", err)
    }
    log.Println("Server exiting")
}

// Stop 停止服务
func (e *Server) Stop() {
    if err := e.server.Shutdown(nil); err != nil {
        panic(err)
    }
}

// Use 中间件
func (e *Server) Use(next gin.HandlerFunc) {
    e.route = e.route.Use(next)
}

// AddRoute 添加路由
func (e *Server) AddRoute(route Route) {
    e.route = e.route.Handle(route.Method, route.Path, route.Handler)
}

// AddRoutes 添加路由
func (e *Server) AddRoutes(routes []Route) {
    for _, route := range routes {
        e.route = e.route.Handle(route.Method, route.Path, route.Handler)
    }
}

// GetEngine 获取路由引擎
func (e *Server) GetEngine() *gin.Engine {
    return e.engine
}

// AddRouteCallback 添加路由
func (e *Server) AddRouteCallback(routerCallback func(*gin.Engine)) {
    routerCallback(e.engine)
}
