package cache

import (
    "fmt"
"sync"

"github.com/gomodule/redigo/redis"
)

// Rediser use redis
type Rediser struct {
    conn    redis.Conn
    subConn redis.Conn
    locker  sync.Mutex
}

var (
    _cache *Rediser

    // ErrNotFound mean can not find the key
    ErrNotFound = fmt.Errorf("the key not found")
)

// NewCache create an cache to read redis
func NewCache(conn redis.Conn) *Rediser {
    if _cache == nil {
        _cache = &Rediser{
            conn: conn,
        }

    }
    return _cache
}

// NewSubCache create an cache support receive message
func NewSubCache(conn, subConn redis.Conn) *Rediser {
    if _cache == nil {
        _cache = &Rediser{
            conn:    conn,
            subConn: subConn,
        }

    }
    subConn.Do("config", "set", "notify-keyspace-events", "KEA")
    return _cache
}

// EventHandler progress message
type EventHandler interface {
    OnMessage(msg redis.Message)
}

// Subscription redis event
func (c *Rediser) Subscription(channel string, handler EventHandler) error {
    psc := redis.PubSubConn{Conn: c.subConn}
    psc.PSubscribe(channel)
    for {
        switch msg := psc.Receive().(type) {
        case redis.Message:
            // fmt.Printf("Message: %s %s\n", msg.Channel, msg.Data)
            handler.OnMessage(msg)
        case redis.Subscription:
            fmt.Printf("Subscription: %s %s %d\n", msg.Kind, msg.Channel, msg.Count)
        case error:
            fmt.Printf("error: %v\n", msg)
        }
    }
}

func (c *Rediser) HSet(key, filed, value interface{}) error {
    return c.doCommand("HSET", key, filed, value)
}

func (c *Rediser) HGet(key, filed interface{}) (string, error) {
    return c.doGetCommand("HGET", key, filed)
}

func (c *Rediser) HGetAll(key interface{}) ([]string, error) {
    return redis.Strings(c.getCommand("HGETALL", key))
}

func (c *Rediser) HGetInt(key, filed string) (int, error) {
    return c.doGetCommandInt("HGET", key, filed)
}

func (c *Rediser) HDel(key, filed string) error {
    return c.doCommand("HDEL", key, filed)
}

func (c *Rediser) Del(key string) error {
    return c.doCommand("DEL", key)
}

func (c *Rediser) Get(key string) (string, error) {
    return c.doGetCommand("GET", key)
}

func (c *Rediser) Expire(key string, ttl int) (error) {
    return c.doCommand("EXPIRE", key, ttl)
}

func (c *Rediser) Set(key string, value interface{}, ttl int) error {
    err := c.doCommand("SET", key, value)
    if err != nil {
        return err
    }
    if ttl != 0 {
        return c.doCommand("EXPIRE", key, ttl)
    }
    return nil
}

func (c *Rediser) getCommand(cmd string, arg ...interface{}) (interface{}, error) {
    c.locker.Lock()
    defer c.locker.Unlock()
    reply, err := c.conn.Do(cmd, arg...)
    if err != nil {
        return nil, err
    }
    if reply == nil {
        return nil, ErrNotFound
    }
    return reply, err
}

func (c *Rediser) doGetCommand(cmd string, arg ...interface{}) (string, error) {
    return redis.String(c.getCommand(cmd, arg...))
}

func (c *Rediser) doGetCommandInt(cmd string, arg ...interface{}) (int, error) {
    return redis.Int(c.getCommand(cmd, arg...))
}

func (c *Rediser) doCommand(cmd string, arg ...interface{}) error {
    c.locker.Lock()
    defer c.locker.Unlock()
    _, err := c.conn.Do(cmd, arg...)
    return err
}
