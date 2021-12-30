package hi

import (
    "fmt"
    "runtime/debug"
    "sync"
    "time"
)

type errorTrack struct {
    err error
    msg string
}

func (et errorTrack) Error() string {
    if et.err != nil {
        return et.err.Error() + "\n" + et.msg
    }
    return et.msg
}

func newErrorTrack(err error, msg string) errorTrack {
    return errorTrack{err: err, msg: msg}
}

type Statistics struct {
    items       []*StatisticsItem
    lock        sync.Mutex
    success     []*StatisticsItem
    fail        []*StatisticsItem
    done        sync.Once
    events      []func(*StatisticsItem)
    parallelism int
    costTotal time.Duration
    costSucc  time.Duration
}

func (s *Statistics) AddEvent(event func(*StatisticsItem)) {
    s.events = append(s.events, event)
}

func (s *Statistics) split() {
    s.done.Do(func() {
        for _, item := range s.items {
            if item.IsSuccess {
                s.success = append(s.success, item)
            } else {
                s.fail = append(s.fail, item)
            }
        }
    })
}

func (s *Statistics) Count() int {
    return len(s.items)
}

func (s *Statistics) CountSucc() int {
    s.split()
    return len(s.success)
}

func (s *Statistics) CountFail() int {
    s.split()
    return len(s.fail)
}

func (s *Statistics) CostTotal() time.Duration {
    if s.costTotal != 0 {
        return s.costTotal
    }
    s.costTotal = 0
    for _, item := range s.items {
        s.costTotal += item.End.Sub(item.Start)
    }
    return s.costTotal
}

func (s *Statistics) CostSucc() time.Duration {
    s.split()
    if s.costSucc != 0 {
        return s.costSucc
    }
    s.costSucc = 0
    for _, item := range s.items {
        s.costSucc += item.End.Sub(item.Start)
    }
    return s.costSucc
}

func (s *Statistics) CostSussAvg() time.Duration {
    return s.CostSucc() / time.Duration(s.CountSucc())
}

func (s *Statistics) Add(item *StatisticsItem) {
    s.lock.Lock()
    s.items = append(s.items, item)
    if s.events != nil {
        for _, event := range s.events {
            event(item)
        }
    }
    s.lock.Unlock()
}

func (s *Statistics) TPS() int {
    return s.parallelism * int(time.Second) / int(s.CostSussAvg())
}

type StatisticsItem struct {
    Start     time.Time
    End       time.Time
    Code      int
    IsSuccess bool
}

// Hi 并发测试
// parallelism 并发数
// n 单协程执行次数
// f 需要执行函数
func Hi(parallelNum, n int, statistics *Statistics, hiFunc func() (int, bool)) (*Statistics, error) {
    if statistics == nil {
        statistics = &Statistics{
            items:       make([]*StatisticsItem, 0, parallelNum*n),
            parallelism: parallelNum,
        }
    }
    wg := sync.WaitGroup{}
    wg.Add(parallelNum)
    var err error
    var lock sync.Mutex
    for i := 0; i < parallelNum; i++ {
        go func() {
            defer func() {
                if gerr := recover(); gerr != nil {
                    lock.Lock()
                    err = newErrorTrack(err, fmt.Sprint(gerr)+"\n"+string(debug.Stack()))
                    lock.Unlock()
                }
            }()
            for j := 0; j < n; j++ {
                item := &StatisticsItem{
                    Start: time.Now(),
                }
                item.Code, item.IsSuccess = hiFunc()
                item.End = time.Now()
                statistics.Add(item)
            }
            wg.Done()
        }()
    }
    wg.Wait()
    return statistics, err
}

// HiDura 并发测试执行指定时间
// parallelism 并发数
// dura 持续时间
// f 需要执行函数
func HiDura(parallelism int, dura time.Duration, hiFunc func() (int, bool)) (*Statistics, error) {
    statistics := &Statistics{
        items:       make([]*StatisticsItem, 0, parallelism*int(int64(dura)/int64(time.Second))),
        parallelism: parallelism,
    }

    chs := make([]chan struct{}, parallelism)
    for i := 0; i < parallelism; i++ {
        chs[i] = make(chan struct{}, 1)
    }

    var err error
    var lock sync.Mutex
    wg := sync.WaitGroup{}
    wg.Add(parallelism)
    for i := 0; i < parallelism; i++ {
        go func(ch chan struct{}) {
            defer func() {
                if gerr := recover(); gerr != nil {
                    lock.Lock()
                    err = newErrorTrack(err, fmt.Sprint(gerr)+"\n"+string(debug.Stack()))
                    lock.Unlock()
                }
            }()
            for {
                select {
                case <-ch:
                    wg.Done()
                    return
                default:
                    item := &StatisticsItem{
                        Start: time.Now(),
                    }
                    item.Code, item.IsSuccess = hiFunc()
                    item.End = time.Now()
                    statistics.Add(item)
                }
            }
        }(chs[i])
    }
    time.AfterFunc(dura, func() {
        for _, item := range chs {
            item <- struct{}{}
        }
    })
    wg.Wait()
    return statistics, err
}
