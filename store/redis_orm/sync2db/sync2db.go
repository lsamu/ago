package sync2db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/weikaishio/distributed_lib/db_lazy"
	"reflect"
	"sync"
	"time"
)

type Sync2DB struct {
	mysqlOrm  *xorm.Engine
	isShowLog bool
	*db_lazy.LazyMysql
	wait *sync.WaitGroup
}

func (s *Sync2DB) IsShowLog(isShow bool) {
	s.isShowLog = isShow
}
func NewSync2DB(mysqlOrm *xorm.Engine, lazyTimeSecond int, wait *sync.WaitGroup) *Sync2DB {
	sync2DB := &Sync2DB{
		mysqlOrm: mysqlOrm,
		wait:     wait,
	}
	sync2DB.LazyMysql = db_lazy.NewLazyMysql(mysqlOrm, lazyTimeSecond)
	go func() {
		go sync2DB.LazyMysql.Exec()
		ListenQuitAndDump()
		sync2DB.LazyMysql.Quit()
		if sync2DB.wait!=nil {
			sync2DB.wait.Done()
		}
	}()
	return sync2DB
}
func (s *Sync2DB) Create2DB(bean interface{}) error {
	err := s.mysqlOrm.Sync(bean)
	if err != nil {
		s.Printfln("mysqlOrm.Sync(%v),err:%v", reflect.TypeOf(bean).Name(), err)
	} else {
		s.Printfln("mysqlOrm.Sync(%v)", reflect.TypeOf(bean).Name())
	}
	return err
}
func (s *Sync2DB) Printfln(format string, a ...interface{}) {
	if s.isShowLog {
		s.Printf(format, a...)
		fmt.Print("\n")
	}
}

func (s *Sync2DB) Printf(format string, a ...interface{}) {
	if s.isShowLog {
		fmt.Printf(fmt.Sprintf("[redis_orm %s] : %s", time.Now().Format("06-01-02 15:04:05"), format), a...)
	}
}