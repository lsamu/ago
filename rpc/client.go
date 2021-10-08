package rpc

import (
    "fmt"
    "google.golang.org/grpc"
)

type (
    Client struct {
        conn *grpc.ClientConn
    }
)

func NewClient() {

}

func (c *Client) Start()  {
    conn, err := grpc.Dial("", grpc.WithInsecure())
    if err != nil {
        fmt.Println(conn)
    }
    c.conn = conn
}

func (c *Client) GetClient() *grpc.ClientConn {
    return c.conn
}

func (c *Client) AddService(regClient func(c *grpc.ClientConn)) {
    regClient(c.conn)
}
