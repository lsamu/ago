package redis_orm

import "github.com/go-redis/redis"

type (
    Conf struct {
        Redis struct {
            Addr     string
            Password string
            DB       int
        }
        MySql struct {
            Host     string
            Database string
            Username string
            Password string
        }
    }

    RedisOrm struct {
        engine *Engine
    }
)

func NewRedisOrm() *RedisOrm {
    redisOrm := &RedisOrm{}
    rediser := redisOrm.redisClient()
    engine := NewEngine(rediser)
    redisOrm.engine = engine
    //_, _ = engine.Schema.ReloadTables()
    //engine.SetSync2DB()
    return redisOrm
}

func (r *RedisOrm) redisClient() (rediser *redis.Client) {
    options := redis.Options{
        Addr:     "127.0.0.1:6379",
        Password: "",
        DB:       1,
    }
    rediser = redis.NewClient(&options)
    return
}

//Sync2DB 同步到数据库
func (r *RedisOrm) Sync2DB() {

}

//Sync2Redis 同步到Redis
func (r *RedisOrm) Sync2Redis() {

}




