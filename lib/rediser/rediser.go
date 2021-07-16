package rediser


import (
    "github.com/gomodule/redigo/redis"
    "github.com/lsamu/ago/lib/cache"
    "time"
)

var redisConn *cache.Rediser

// InitRedis 连接redis
func InitRedis(address string,password string,db int) {
    conn, err := redis.Dial("tcp", address,
        redis.DialReadTimeout(10*time.Second),
        redis.DialWriteTimeout(10*time.Second),
        redis.DialConnectTimeout(10*time.Second),
        redis.DialDatabase(db),
        redis.DialPassword(password),
    )
    if err != nil {
        panic(err)
    }
    redisConn = cache.NewSubCache(conn, conn)
    return
}

// RedisConn 连接redis
func RedisConn() *cache.Rediser {
    return redisConn
}

