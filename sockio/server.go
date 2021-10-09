package sockio

import (
    "fmt"
    "github.com/gin-gonic/gin"
    socketio "github.com/googollee/go-socket.io"
    "log"
)

type Server struct {
    conf     SockConf
    engine   *gin.Engine
    socketIO *socketio.Server
}

//SocketEvent SocketEvent
type SocketEvent struct {
    NameSpace string
    Event     string
    Handler   interface{} //func(s socketio.Conn, msg string),func(s socketio.Conn, msg string) string,func(s socketio.Conn) string,
}

func NewServer(conf SockConf) *Server {
    engine := gin.New()
    server := socketio.NewServer(nil)
    server.OnConnect("/", func(s socketio.Conn) error {
        s.SetContext("")
        log.Println("connected:", s.ID())
        return nil
    })
    server.OnError("/", func(s socketio.Conn, err error) {
        log.Println("error:", err)
    })
    server.OnDisconnect("/", func(s socketio.Conn, reason string) {
        log.Println("closed", reason)
    })
    return &Server{
        engine:   engine,
        socketIO: server,
        conf:     conf,
    }
}

func (e *Server) Start() {
    go func() {
        if err := e.socketIO.Serve(); err != nil {
            log.Fatalf("socketio listen error: %s\n", err)
        }
    }()
    defer e.socketIO.Close()
    e.engine.GET("/socket.io/*any", gin.WrapH(e.socketIO))
    e.engine.POST("/socket.io/*any", gin.WrapH(e.socketIO))
    //e.engine.StaticFS("/public", http.Dir("../asset"))
    bind := fmt.Sprintf("%s:%d", e.conf.Host, e.conf.Port)
    if err := e.engine.Run(bind); err != nil {
        log.Fatal("failed run app: ", err)
    }
}

func (e *Server) Stop() {
    e.socketIO.Close()
}

// Use 中间件
func (e *Server) Use(next gin.HandlerFunc) {
    e.engine.Use(next)
}

// GetSocketIO GetSocketIO
func (e *Server) GetSocketIO() *socketio.Server {
    return e.socketIO
}

// AddEvent 添加事件
func (e *Server) AddEvent(event SocketEvent) {
    if event.NameSpace == "" {
        event.NameSpace = "/"
    }
    e.socketIO.OnEvent(event.NameSpace, event.Event, event.Handler)
}

// AddEvents 添加多个事件
func (e *Server) AddEvents(events []SocketEvent) {
    for _, event := range events {
        if event.NameSpace == "" {
            event.NameSpace = "/"
        }
        e.socketIO.OnEvent(event.NameSpace, event.Event, event.Handler)
    }
}
