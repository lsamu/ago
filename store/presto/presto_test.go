package presto

import (
    "fmt"
    "testing"
)

var presto *Presto

func TestMain(m *testing.M) {
    db, err := Open("http://admin@10.30.32.226:8523?catalog=hive&schema=hub")
    if err != nil {
        fmt.Println(err)
    }
    presto = db
    m.Run()
}

func TestOpen(t *testing.T) {
    var aa struct {
        ss int
    }
    err := presto.Query("select * from hive.hub.dwi_nginx_log where pt=20200505 limit 1").First(&aa).Error
    if err != nil {
        fmt.Println(err)
    }
}
