// +build !windows

package sync2db

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"syscall"
)

type SignalHandleFunc func(sig os.Signal, reload SignalRelodFunc) (ret bool)
type SignalRelodFunc func()

// ListenSignal 监听信号 signals, 当收到其中一个信号时调用 handler
//函数将会阻塞直到指定的信号到来且 handler 处理信号后返回true(如果返回false,会继续接受信号)
func ListenSignal(handler SignalHandleFunc, reload SignalRelodFunc, signals ...os.Signal) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, signals...)
	for {
		sig := <-sigChan
		if handler(sig, reload) {
			break
		}
	}
}

//ListenQuitAndDump 函数将会阻塞直到 INT/USR1/USR2 信号到来
func ListenQuitAndDump() {
	ListenSignal(QuitAndDumpAndReload, nil, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)
}

func ListenQuitAndDumpAndReload(reload SignalRelodFunc) {
	ListenSignal(QuitAndDumpAndReload, reload, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2)
}

// QuitAndDumpAndReload 这是一个 SignalHandleFunc,用于退出或dump进程
// 退出监听: syscall.SIGINT, syscall.SIGUSR1
// dump监听: syscall.SIGUSR2
// 使用 kill 命令时可以带上信号参数:
//
//	kill -s INT <pid> 杀进程
//	kill -s USR1 <pid> reload配置
//	kill -s USR2 <pid> dump内存堆栈
func QuitAndDumpAndReload(sig os.Signal, reload SignalRelodFunc) bool {
	switch sig {
	case syscall.SIGINT:
		return true
	case syscall.SIGUSR1:
		if reload != nil {
			reload()
		}
	case syscall.SIGUSR2:
		filename := filepath.Base(os.Args[0]) + ".dump"
		dumpOut, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
		if err == nil {
			for _, name := range []string{"goroutine", "heap", "block"} {
				p := pprof.Lookup(name)
				if p == nil {
					continue
				}
				name = strings.ToUpper(name)
				fmt.Fprintf(dumpOut, "-----BEGIN %s-----\n", name)
				p.WriteTo(dumpOut, 2)
				fmt.Fprintf(dumpOut, "\n-----END %s-----\n", name)
			}
			dumpOut.Close()
		}
	}
	return false
}
