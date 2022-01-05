package redis_orm

import (
    "github.com/go-redis/redis"
)

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
    // _, _ = engine.Schema.ReloadTables()
    // engine.SetSync2DB()
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

// Sync2DB 同步到数据库
func (r *RedisOrm) Sync2DB() {

}

// Sync2Redis 同步到Redis
func (r *RedisOrm) Sync2Redis() {

}

func Open(conf *Conf) (redisRom *RedisOrm) {
    return  NewRedisOrm()
}

// Model Model
func (r *RedisOrm) Model(value interface{}) (redisRom *RedisOrm) {
    return
}

// Table Table
func (r *RedisOrm) Table(name string) (redisRom *RedisOrm) {
    return
}

// Where Where
func (r *RedisOrm) Where(query interface{}, args ...interface{}) (redisRom *RedisOrm) {
    return
}

// Or Or
func (r *RedisOrm) Or(query interface{}, args ...interface{}) (redisRom *RedisOrm) {
    return
}

// Not Not
func (r *RedisOrm) Not(query interface{}, args ...interface{}) (redisRom *RedisOrm) {
    return
}

// Limit Limit
func (r *RedisOrm) Limit(limit interface{}) (redisRom *RedisOrm) {
    return
}

// Offset Offset
func (r *RedisOrm) Offset(offset interface{}) (redisRom *RedisOrm) {
    return
}

// Select Select
func (r *RedisOrm) Select(query interface{}, args ...interface{}) () {
    return
}

// Find Find
func (r *RedisOrm) Find(out interface{}, where ...interface{}) (redisRom *RedisOrm) {
    // searchCon := NewSearchConditionV2(faq.Unique, faq.Unique, 111)
    // var ary []models.Faq
    // r.engine.Find(0,100,searchCon,&ary)
    return
}

// Count Count
func (r *RedisOrm) Count(value interface{}) (redisRom *RedisOrm) {
    return
}

// First First
func (r *RedisOrm) First(out interface{}, where ...interface{}) (redisRom *RedisOrm) {
    r.engine.Get(out)
    return
}

// Create Create
func (r *RedisOrm) Create(value interface{}) (redisRom *RedisOrm) {
    return
}

// Update Update
func (r *RedisOrm) Update(values interface{}, ignoreProtectedAttrs ...bool) (redisRom *RedisOrm) {
    return
}

// Delete Delete
func (r *RedisOrm) Delete(value interface{}, where ...interface{}) (redisRom *RedisOrm) {
    return
}

// Group Group
func (r *RedisOrm) Group(query string) (redisRom *RedisOrm) {
    return
}

// Order Order
func (r *RedisOrm) Order(value interface{}, reorder ...bool) (redisRom *RedisOrm) {
    return
}

// Raw Raw
func (r *RedisOrm) Raw(sql string, values ...interface{}) (redisRom *RedisOrm) {
    return
}

// Exec Exec
func (r *RedisOrm) Exec(sql string, values ...interface{}) (redisRom *RedisOrm) {
    return
}
