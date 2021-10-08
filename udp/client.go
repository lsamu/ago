package udp

import (
    "fmt"
    "net"
)

type Client struct {

}

func NewClient() {

}

func (c Client) Start()  {
    conn, err := net.Dial("udp", "127.0.0.1:8002")
    if err != nil {
        fmt.Println("net.Dial err:", err)
        return
    }
    defer conn.Close()

    conn.Write([]byte("Hello! I'm client in UDP!"))

    buf := make([]byte, 1024)
    n, err1 := conn.Read(buf)
    if err1 != nil {
        return
    }
    fmt.Println("服务器发来：", string(buf[:n]))
}