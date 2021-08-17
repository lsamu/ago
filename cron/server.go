package cron

import (
    cron "github.com/robfig/cron/v3"
    "log"
    "os"
    "os/signal"
)

type Server struct {
    cron *cron.Cron
}

func NewServer() *Server {
    return &Server{
        cron: cron.New(cron.WithSeconds()),
    }
}

func (s *Server) Start() {
    s.cron.Run()
    //t1 := time.NewTimer(time.Second * 10)
    //for {
    //    select {
    //    case <-t1.C:
    //        t1.Reset(time.Second * 10)
    //    }
    //}
    quit := make(chan os.Signal)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutdown Cron ...")
}

func (s *Server) Stop() {
    s.cron.Stop()
}

func (s *Server) Use() {

}

func (s *Server) AddFunc(spec string, cmd func()) {
    addFunc, err := s.cron.AddFunc(spec, cmd)
    if err != nil {
        log.Fatalf("%+v", err)
        return
    }
    log.Printf("task start.num:%+v", addFunc)
}
