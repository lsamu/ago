package sock

import (
    "github.com/gorilla/websocket"
    "log"
    "net/url"
)

func NewClient() {

}

func client()  {
    u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
    log.Printf("connecting to %s", u.String())
    c, _, err := websocket.DefaultDialer.Dial(u.String(),nil)
    if err != nil {
        log.Fatal("dial:", err)
    }
    defer c.Close()
    go func() {
        log.Println("start message loop")
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("read:", err)
                return
            }
            log.Printf("recv: %s", message)
        }
    }()
}
