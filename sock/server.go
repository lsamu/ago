package sock

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "io"
    "log"
    "net/http"
)

/*
var upgrader = websocket.Upgrader{
    ReadBufferSize:    4096,
    WriteBufferSize:   4096,
    EnableCompression: true,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}*/

var upgrader = websocket.Upgrader{
    // 解决跨域问题
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

//原生
func ws(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer conn.Close()
    for {
        mt, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        err = conn.WriteMessage(mt, message)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
}

//gin
func gins(c *gin.Context)  {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer conn.Close()
    for {
        mt, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }
        log.Printf("recv: %s", message)
        err = conn.WriteMessage(mt, message)
        if err != nil {
            log.Println("write:", err)
            break
        }
    }
}

func NextRead(conn *websocket.Conn)  {
    for {
        mt, r, err := conn.NextReader()
        if err != nil {
            if err != io.EOF {
                log.Println("NextReader:", err)
            }
            return
        }
        if mt == websocket.TextMessage {
            //r = &validator{r: r}
            fmt.Println(r)
        }
    }
}