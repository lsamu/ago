package logger

import (
    "encoding/json"
"os"

logging "github.com/op/go-logging"
)

// DefaultWebFormatter for gin web log
var DefaultWebFormatter = logging.MustStringFormatter(
    "%{color}[%{level:.3s}] %{time:2006/01/02 - 15:04:05} | %{shortfile:s} | %{shortfunc:s} | %{id:06x}%{color:reset} %{message}",
)

// NewLogger with category, formatter, extraCalldepth
func NewLogger(category string, formatter logging.Formatter, extraCalldepth int) *logging.Logger {
    log := logging.MustGetLogger(category)
    // For demo purposes, create two backend for os.Stderr.
    backendError := logging.NewLogBackend(os.Stderr, "", 0)
    backendAll := logging.NewLogBackend(os.Stdout, "", 0)
    backend2Formatter := logging.NewBackendFormatter(backendAll, formatter)
    backend1Leveled := logging.AddModuleLevel(backendError)
    backend1Leveled.SetLevel(logging.ERROR, "")
    // Set the backends to be used.
    logging.SetBackend(backend1Leveled, backend2Formatter)
    log.ExtraCalldepth = extraCalldepth
    return log
}

// Logger are dispard
var Logger = NewLogger("", DefaultWebFormatter, 0)

var _PriLogger = NewLogger("", DefaultWebFormatter, 1)

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func Fatal(args ...interface{}) {
    _PriLogger.Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
    _PriLogger.Fatalf(format, args...)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func Panic(args ...interface{}) {
    _PriLogger.Panic(args...)
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func Panicf(format string, args ...interface{}) {
    _PriLogger.Panicf(format, args...)
}

// Critical logs a message using CRITICAL as log level.
func Critical(args ...interface{}) {
    _PriLogger.Critical(args...)
}

// Criticalf logs a message using CRITICAL as log level.
func Criticalf(format string, args ...interface{}) {
    _PriLogger.Criticalf(format, args...)
}

// Error logs a message using ERROR as log level.
func Error(args ...interface{}) {
    _PriLogger.Error(args...)
}

// Errorf logs a message using ERROR as log level.
func Errorf(format string, args ...interface{}) {
    _PriLogger.Errorf(format, args...)
}

// Warning logs a message using WARNING as log level.
func Warning(args ...interface{}) {
    _PriLogger.Warning(args...)
}

// Warningf logs a message using WARNING as log level.
func Warningf(format string, args ...interface{}) {
    _PriLogger.Warningf(format, args...)
}

// Notice logs a message using NOTICE as log level.
func Notice(args ...interface{}) {
    _PriLogger.Notice(args...)
}

// Noticef logs a message using NOTICE as log level.
func Noticef(format string, args ...interface{}) {
    _PriLogger.Noticef(format, args...)
}

// Info logs a message using INFO as log level.
func Info(args ...interface{}) {
    _PriLogger.Info(args...)
}

// Infof logs a message using INFO as log level.
func Infof(format string, args ...interface{}) {
    _PriLogger.Infof(format, args...)
}

// Debug logs a message using DEBUG as log level.
func Debug(args ...interface{}) {
    _PriLogger.Debug(args...)
}

// Debugf logs a message using DEBUG as log level.
func Debugf(format string, args ...interface{}) {
    _PriLogger.Debugf(format, args...)
}

// Dump data
func Dump(v ...interface{}) {
    for _, val := range v {
        data, _ := json.Marshal(val)
        _PriLogger.Debug(string(data))
    }
}
