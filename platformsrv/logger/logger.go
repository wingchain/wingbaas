
package logger

import (
  "fmt"
  "os"
  "runtime"
  "github.com/op/go-logging"
)

var log = logging.MustGetLogger("wingbaas")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
  `%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, format)
    logging.SetBackend(backendFormatter).SetLevel(logging.DEBUG, "")
}

func SetLogInit() {
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, format)
    logging.SetBackend(backendFormatter).SetLevel(logging.DEBUG, "")
}

func getfileName(file string) string {
    short := file
    for i := len(file) - 1; i > 0; i-- {
      if file[i] == '/' {
        short = file[i+1:]
        break
      }
    }
    return short
}

func Debug(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    str := fmt.Sprintf("%s %d ", getfileName(file), line)
    var newArgs []interface{}
    newArgs = append(newArgs,str)
    newArgs = append(newArgs,args...)
    log.Debug(newArgs...)
}

func Debugf(format string, args ...interface{}) {
     _, file, line, _ := runtime.Caller(1)
    newFormatStr := fmt.Sprintf("%s %d ", getfileName(file), line) + format
    log.Debugf(newFormatStr, args...)
}

func Info(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    str := fmt.Sprintf("%s %d ", getfileName(file), line)
    var newArgs []interface{}
    newArgs = append(newArgs,str)
    newArgs = append(newArgs,args...)
    log.Info(newArgs...)
}

func Infof(format string, args ...interface{}) {
     _, file, line, _ := runtime.Caller(1)
    newFormatStr := fmt.Sprintf("%s %d ", getfileName(file), line) + format
    log.Infof(newFormatStr, args...)
}

func Warning(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    str := fmt.Sprintf("%s %d ", getfileName(file), line)
    var newArgs []interface{}
    newArgs = append(newArgs,str)
    newArgs = append(newArgs,args...)
    log.Warning(newArgs...)
}

func Warningf(format string, args ...interface{}) {
     _, file, line, _ := runtime.Caller(1)
    newFormatStr := fmt.Sprintf("%s %d ", getfileName(file), line) + format
    log.Warningf(newFormatStr, args...)
}

func Error(args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    str := fmt.Sprintf("%s %d ", getfileName(file), line)
    var newArgs []interface{}
    newArgs = append(newArgs,str)
    newArgs = append(newArgs,args...)
    log.Error(newArgs...)
}

func Errorf(format string, args ...interface{}) {
     _, file, line, _ := runtime.Caller(1)
    newFormatStr := fmt.Sprintf("%s %d ", getfileName(file), line) + format
    log.Errorf(newFormatStr, args...)
}
