package defense

import (
    "context"
    "errors"
    //"github.com/go-redis/redis"
    "github.com/go-redis/redis/v8"
    "sync"
    "time"
)

// 基于redis的防暴力接口尝试防御简单实现
// 示例：
// sDefense := defense.New(redis, 10 * time.Minute, 30)
//
// 举例1：某个账号登录错误次数默认设置防御--10分钟内错误次数不得超过30次
// if sDefense.Defense("login_account") != nil {
//		return errors.New("密码错误次数过多，请稍后再试或尝试找回密码")
// }
//
// 举例2：某个ip登录错误次数默认设置防御--10分钟内错误次数不得超过30次
// if sDefense.Defense("127.0.0.1") != nil {
//		return errors.New("密码错误次数过多，请稍后再试或尝试找回密码")
// }
//
// 举例3：某个ip登录错误次数自定义设置防御--30分钟内错误次数不得超过20次
// if sDefense.DefenseCustom("127.0.0.1", 30 * time.Minute, 20) != nil {
//		return errors.New("密码错误次数过多，请稍后再试或尝试找回密码")
// }
//
// 举例4：登录成功需要释放掉某个防御
// sDefense.Release("login_account")
// sDefense.Release("127.0.0.1")

// ErrInDefense 定义处于防御模式需要拦截的错误
var ErrInDefense = errors.New("error of in defense")

// SimpleDefense 简单防御实现结构
type SimpleDefense struct {
    redis    *redis.Client
    duration time.Duration
    times    int64
    mu       sync.Mutex
}

// New 创建一个简单实现的防暴力破解实例
// @param redis go-redis/redis v7 对象示例
// @param defenseDuration 默认防御间隔时长设置，譬如：1分钟内最大尝试次数不得超过5次，此处传值 1 * time.Minute
// @param defenseTimes    默认防御间隔次数设置，譬如：1分钟内最大尝试次数不得超过5次，此处传值 5
func New(redis *redis.Client, defenseDuration time.Duration, defenseTimes int64) *SimpleDefense {
    return &SimpleDefense{
        redis:    redis,
        duration: defenseDuration,
        times:    defenseTimes,
        mu:       sync.Mutex{},
    }
}

// Defense 设置防御，未超过初始化条件设置防御成功返回nil，超过初始化条件设置防御失败返回error即触发了防御条件需要拦截
// @param defenseKey 按默认初始化策略检测指定key防御
func (s *SimpleDefense) Defense(defenseKey string) error {
    return s.DefenseCustom(defenseKey, s.duration, s.times)
}

// DefenseCustom 设置自定义防御，未超过初始化条件设置防御成功返回nil，超过初始化条件设置防御失败返回error即触发了防御条件需要拦截
// @param defenseKey 	  按自定义策略检测指定key防御
// @param defenseDuration 防御策略有效期
// @param defenseTimes    防御策略有效期内的最大次数
func (s *SimpleDefense) DefenseCustom(defenseKey string, defenseDuration time.Duration, defenseTimes int64) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    ctx := context.TODO()
    // 没有这个key，设置key并返回
    res := s.redis.Get(ctx, defenseKey)
    if res.Err() == redis.Nil {
        s.redis.Set(ctx, defenseKey, 1, defenseDuration)
        return nil
    }

    // 尝试读取转换已有次数出错，返回拦截状态
    tryTimes, err := res.Int64()
    if err != nil {
        return ErrInDefense
    }

    // 防御次数超标：返回防御拦截同时再次延长key有效期
    if tryTimes+1 >= defenseTimes {
        _ = s.redis.Expire(ctx, defenseKey, defenseDuration)
        return ErrInDefense
    }

    // 不处于拦截状态，累加1
    _ = s.redis.Set(ctx, defenseKey, tryTimes+1, defenseDuration)

    return nil
}

// Release 释放指定防御
func (s *SimpleDefense) Release(defenseKey string) {
    _ = s.redis.Del(context.TODO(), defenseKey)
}