package expired_map

import (
    "sync"
    "time"
)

const delChannelCap = 100

type (
    val struct {
        data        interface{}
        expiredTime int64
    }

    ExpiredMap struct {
        m       map[interface{}]*val
        lck     *sync.Mutex
        timeMap map[int64][]interface{}
        stop    chan struct{}
    }
    delMsg struct {
        keys []interface{}
        t    int64
    }
)

func NewExpiredMap() *ExpiredMap {
    e := ExpiredMap{
        m:       make(map[interface{}]*val),
        lck:     new(sync.Mutex),
        timeMap: make(map[int64][]interface{}),
        stop:    make(chan struct{}),
    }
    go e.run(time.Now().Unix())
    return &e
}

func (e *ExpiredMap) run(now int64) {
    t := time.NewTicker(time.Second * 1)
    defer t.Stop()

    delCh := make(chan *delMsg, delChannelCap)
    go func() {
        for v := range delCh {
            e.multiDelete(v.keys, v.t)
        }
    }()

    for {
        select {
        case <-t.C:
            now++ //这里用now++的形式，直接用time.Now().Unix()可能会导致时间跳过1s，导致key未删除。
            e.lck.Lock()
            if keys, found := e.timeMap[now]; found {
                e.lck.Unlock()
                delCh <- &delMsg{keys: keys, t: now}
            } else {
                e.lck.Unlock()
            }
        case <-e.stop:
            close(delCh)
            return
        }
    }
}

func (e *ExpiredMap) Set(key, value interface{}, expireSeconds int64) {
    if expireSeconds <= 0 {
        return
    }
    e.lck.Lock()
    defer e.lck.Unlock()
    expiredTime := time.Now().Unix() + expireSeconds
    e.m[key] = &val{
        data:        value,
        expiredTime: expiredTime,
    }
    e.timeMap[expiredTime] = append(e.timeMap[expiredTime], key)
}

func (e *ExpiredMap) Get(key interface{}) (found bool, value interface{}) {
    e.lck.Lock()
    defer e.lck.Unlock()
    if found = e.checkDeleteKey(key); !found {
        return
    }
    value = e.m[key].data
    return
}

func (e *ExpiredMap) Delete(key interface{}) {
    e.lck.Lock()
    delete(e.m, key)
    e.lck.Unlock()
}

func (e *ExpiredMap) Remove(key interface{}) {
    e.Delete(key)
}

func (e *ExpiredMap) Length() int {
    e.lck.Lock()
    defer e.lck.Unlock()
    return len(e.m)
}

func (e *ExpiredMap) Size() int {
    return e.Length()
}

func (e *ExpiredMap) TTL(key interface{}) int64 {
    e.lck.Lock()
    defer e.lck.Unlock()
    if !e.checkDeleteKey(key) {
        return -1
    }
    return e.m[key].expiredTime - time.Now().Unix()
}

func (e *ExpiredMap) Clear() {
    e.lck.Lock()
    defer e.lck.Unlock()
    e.m = make(map[interface{}]*val)
    e.timeMap = make(map[int64][]interface{})
}

func (e *ExpiredMap) checkDeleteKey(key interface{}) bool {
    if val, found := e.m[key]; found {
        if val.expiredTime <= time.Now().Unix() {
            delete(e.m, key)
            return false
        }
        return true
    }
    return false
}

func (e *ExpiredMap) multiDelete(keys []interface{}, t int64) {
    e.lck.Lock()
    defer e.lck.Unlock()
    delete(e.timeMap, t)
    for _, key := range keys {
        delete(e.m, key)
    }
}