package tcp

import (
    "fmt"
    "net"
)

type Server struct {

}

func NewServer() *Server{
    return &Server{}
}

func (s *Server) Start() {
    fmt.Println("Starting the server ...")
    listener, err := net.Listen("tcp", "localhost:50000")
    if err != nil {
        fmt.Println("Error listening", err.Error())
        return //终止程序
    }
    // 监听并接受来自客户端的连接
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting", err.Error())
            return // 终止程序
        }
        go s.doServerStuff(conn)
    }
}
func (s *Server) doServerStuff(conn net.Conn) {
    for {
        buf := make([]byte, 512)
        len, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading", err.Error())
            return //终止程序
        }
        fmt.Printf("Received data: %v\n", string(buf[:len]))
    }
}