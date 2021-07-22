package main

import (
    "fmt"
    "github.com/robfig/cron/v3"
    "time"
)

func main() {
    //每隔5秒执行一次：*/5 * * * * ?
    //
    //每隔1分钟执行一次：0 */1 * * * ?
    //
    //每天23点执行一次：0 0 23 * * ?
    //
    //每天凌晨1点执行一次：0 0 1 * * ?
    //
    //每月1号凌晨1点执行一次：0 0 1 1 * ?
    //
    //每月最后一天23点执行一次：0 0 23 L * ?
    //
    //每周星期天凌晨1点实行一次：0 0 1 ? * L
    //
    //在26分、29分、33分执行一次：0 26,29,33 * * * ?
    //
    //每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?

    c:=cron.New()
    _, _ = c.AddFunc("*/5 * * * * *", func() {
        fmt.Println(time.Now())
    })
    defer c.Stop()
    c.Start()

    // 这是一个使用time包实现的定时器，与cron做对比
    t1 := time.NewTimer(time.Second * 10)
    for {
        select {
        case <-t1.C:
            t1.Reset(time.Second * 10)
        }
    }
    //quit := make(chan os.Signal)
    //signal.Notify(quit, os.Interrupt)
    //<-quit
    //log.Println("Shutdown Server ...")
}
