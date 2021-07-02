package queue


import (
    "container/list"
    "errors"
    "fmt"
    "sync"
)

// Queue 队列
type Queue struct {
    List  *list.List
    Mutex *sync.RWMutex
}

// GetQueue 获取队列
func GetQueue() *Queue {
    return &Queue{
        List:  list.New(),
        Mutex: new(sync.RWMutex),
    }
}

// Push 入队列
func (queue *Queue) Push(data interface{}) {
    if data == nil {
        return
    }
    queue.Mutex.Lock()
    defer queue.Mutex.Unlock()
    queue.List.PushBack(data)
}

// Pop 出队列
func (queue *Queue) Pop() (interface{}, error) {
    queue.Mutex.Lock()
    defer queue.Mutex.Unlock()
    if element := queue.List.Front(); element != nil {
        queue.List.Remove(element)
        return element.Value, nil
    }
    return nil, errors.New("pop failed")
}

// Clear 清除队列
func (queue *Queue) Clear() {
    queue.Mutex.Lock()
    defer queue.Mutex.Unlock()
    for element := queue.List.Front(); element != nil; {
        elementNext := element.Next()
        queue.List.Remove(element)
        element = elementNext
    }
}

// Len 获取队列长度
func (queue *Queue) Len() int {
    queue.Mutex.RLock()
    defer queue.Mutex.RUnlock()
    return queue.List.Len()
}

// Show 打印队列
func (queue *Queue) Show() {
    queue.Mutex.RLock()
    defer queue.Mutex.RUnlock()
    for item := queue.List.Front(); item != nil; item = item.Next() {
        fmt.Println(item.Value)
    }
}
