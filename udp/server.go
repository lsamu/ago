package udp

import (
    "fmt"
    "net"
)

type Server struct {

}

func NewServer() {

}

func (s Server) Start()  {
    //创建监听的地址，并且指定udp协议
    udp_addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8002")
    if err != nil {
        fmt.Println("ResolveUDPAddr err:", err)
        return
    }
    conn, err := net.ListenUDP("udp", udp_addr)    //创建数据通信socket
    if err != nil {
        fmt.Println("ListenUDP err:", err)
        return
    }
    defer conn.Close()

    buf := make([]byte, 1024)
    n, raddr, err := conn.ReadFromUDP(buf)        //接收客户端发送过来的数据，填充到切片buf中。
    if err != nil {
        return
    }
    fmt.Println("客户端发送：", string(buf[:n]))

    _, err = conn.WriteToUDP([]byte("nice to see u in udp"), raddr) // 向客户端发送数据
    if err != nil {
        fmt.Println("WriteToUDP err:", err)
        return
    }
}