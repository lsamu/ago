package redis_pubsub

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type RdsPubSubMsg struct {
	redisClient redis.Cmdable
	subCli      *redis.PubSub
	updateLock  sync.RWMutex
	newMsgRev   map[string]func(msg interface{})
	lock        sync.RWMutex
	running     int32
	newChannel  chan struct{}
}

var (
	rdsPubSubMsgClient *RdsPubSubMsg
	onceMsg            sync.Once
	ERR_RedisNotInit   = errors.New("redis not init")
)

func SharedRdsSubscribMsgInstance() *RdsPubSubMsg {
	onceMsg.Do(func() {
		rdsPubSubMsgClient = &RdsPubSubMsg{
			newMsgRev:  make(map[string]func(msg interface{})),
			newChannel: make(chan struct{}),
		}
	})
	return rdsPubSubMsgClient
}

func (r *RdsPubSubMsg) AddSubscribe(channel string, onRevMsg func(msg interface{})) {
	r.updateLock.Lock()
	defer r.updateLock.Unlock()
	r.newMsgRev[channel] = onRevMsg
	go func() {
		r.newChannel <- struct{}{}
	}()
}

func (r *RdsPubSubMsg) Publish(channel string, msg interface{}) error {
	//if r.redisClient != nil {
		_, err := r.redisClient.Publish(channel, msg).Result()
		return err
	//} else {
	//	return ERR_RedisNotInit
	//}
}

func (r *RdsPubSubMsg) Quit() {
	log.Println("RdsPubSub ready quit")
	atomic.SwapInt32(&r.running, 0)
	log.Println("RdsPubSubquit ok")
}
func (r *RdsPubSubMsg) IsRunning() bool {
	return atomic.LoadInt32(&r.running) != 0
}
func (r *RdsPubSubMsg) Set(rdsCli redis.Cmdable) {
		r.redisClient = rdsCli
}
func (r *RdsPubSubMsg) StartSubscription() {
	if !atomic.CompareAndSwapInt32(&r.running, 0, 1) {
		return
	}

	log.Println("StartSubscription")
	defer atomic.CompareAndSwapInt32(&r.running, 1, 0)

	sleepSecond := 3
	for r.IsRunning() {
		select {
		case <-time.After(1 * time.Second):
			if r.redisClient == nil {
				log.Println("StartSubscription rdsPubSubClient.redisClient is nil")
				break
			}
			for r.IsRunning() {
				var channels []string
				r.updateLock.RLock()
				for k := range r.newMsgRev {
					channels = append(channels, k)
				}
				r.updateLock.RUnlock()
				var subCli *redis.PubSub
				switch typ := r.redisClient.(type) {
				case *redis.ClusterClient:
					subCli = typ.PSubscribe(channels...)
				case *redis.Client:
					subCli = typ.PSubscribe(channels...)
				default:
					log.Fatal("invalid redisClient:%v", r.redisClient)
					return
				}
				isOpen, _ := r.subscription(subCli, sleepSecond, channels)

				if !isOpen {
					log.Println("StartSubscription sub.Channel() isClose, u.renewSubClient")
					time.Sleep(time.Duration(sleepSecond) * time.Second)
					sleepSecond += sleepSecond
				}else{
					_ = subCli.Close()
				}
			}
		}
	}
	log.Println("quit StartSubscription")
}
func (r *RdsPubSubMsg) subscription(subCli *redis.PubSub, sleepSecond int, channels []string) (isOpen bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 1<<16)
			buf = buf[:runtime.Stack(buf, true)]
			switch typ := e.(type) {
			case error:
				err = typ
			case string:
				err = errors.New(typ)
			default:
				err = fmt.Errorf("%v", typ)
			}
			log.Printf("==== STACK TRACE BEGIN ====\npanic: %v\n%s\n===== STACK TRACE END =====", err, string(buf))
		}
	}()
	for {
		select {
		case msg, isOpen := <-subCli.Channel():
			if isOpen {
				r.updateLock.RLock()
				onRevMsg, has := r.newMsgRev[msg.Channel]
				r.updateLock.RUnlock()
				if has {
					onRevMsg(msg)
				} else {
					log.Printf("r.newMsgRev[%s] !has", msg.Channel)
				}
			} else {
				return false, nil
			}
		case <-r.newChannel:
			log.Println("StartSubscription new channel rev")
			return true, nil
		case <-time.After(30 * time.Minute):
			log.Println("StartSubscription <-sub.Channel() timeout for 30 min")
			return true, nil
		}
	}
}
