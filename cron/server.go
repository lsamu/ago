package cron

import cron "github.com/robfig/cron/v3"

type Server struct {
    cron *cron.Cron
}

func NewServer() *Server {
    return &Server{
        cron: cron.New(),
    }
}

func (s *Server) Start() {
    s.cron.Start()
    defer s.cron.Stop()
}

func (s *Server) Use() {

}

func (s *Server) AddFun(spec string, cmd func()) (cron.EntryID, error) {
    return s.cron.AddFunc(spec, cmd)
}
