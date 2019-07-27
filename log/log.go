package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	syslog "log"
	"os"
	"path"
	"regexp"
	"sync"
	"time"
)

/**
* 日志级别
 */
type Level int

const (
	FATAL Level = iota
	ERROR
	WARNING
	INFO
	DEBUG
	ALL
)

// 扩展系统logger
// 实现trace_id跟踪
// 实现日志切割/删除

type Log struct {
	logger     *syslog.Logger
	config     LogConfig
	traceId    string // 请求id，用于调用链跟踪
	rotateTime string // .log -> .log.(rotateTime)
	mu         sync.Mutex
}

type LogConfig struct {
	Type          string // syslog/stderr/std/file
	Level         string // DEBUG/INFO/WARNING/ERROR/FATAL
	Dir           string // 文件目录
	FileName      string // 文件名
	RotateByHour  bool   // 按小时切割
	RotateByDaily bool   // 按天切割
	KeepDays      int    // 保留天数
}

var l Log

func Init(path string) error {
	var conf LogConfig
	if res, err := ioutil.ReadFile(path); err != nil {
		return errors.New("error opening conf file=" + path)
	} else {
		if err := json.Unmarshal(res, &conf); err != nil {
			msg := fmt.Sprintf("error parsing conf file=%s, err=%s", path, err.Error())
			return errors.New(msg)
		}
	}

	SetConfig(conf)
	SetLogger(newLogger(conf))
	SetRotateTime(getCurrentTime(conf))

	// 日志切割
	go rotateDaemon()

	return nil
}

func SetTraceId(traceId string) {
	l.traceId = traceId
}

func SetConfig(conf LogConfig) {
	l.config = conf
}

func SetLogger(logger *syslog.Logger) {
	l.logger = logger
}

func SetRotateTime(rotateTime string) {
	l.rotateTime = rotateTime
}

// 以下日志输出函数
// 规定所有非格式化输出参数为map，方便合并trace_id以及

func Debug(args map[string]interface{}) {
	print(DEBUG, args)
}

func Debugf(format string, args ...interface{}) {
	printf(DEBUG, format, args...)
}

func Info(args map[string]interface{}) {
	print(INFO, args)
}

func Infof(format string, args ...interface{}) {
	printf(INFO, format, args...)
}

func Warning(args map[string]interface{}) {
	print(WARNING, args)
}

func Warningf(format string, args ...interface{}) {
	printf(WARNING, format, args...)
}

func Error(args map[string]interface{}) {
	print(ERROR, args)
}

func Errorf(format string, args ...interface{}) {
	printf(ERROR, format, args...)
}

func Fatal(args map[string]interface{}) {
	print(FATAL, args)
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	printf(FATAL, format, args...)
	os.Exit(1)
}

/*
* 非格式化输出，合并trace_id
 */
func print(level Level, m map[string]interface{}) {
	if level > stringToLevel(l.config.Level) {
		return
	}
	if l.traceId != "" {
		m["trace_id"] = l.traceId
	}
	l.logger.SetPrefix(getPrefixByLevel(level))

	header := header()
	body := mapToStr(m)

	buf := contentToBuffer(header, body)
	l.logger.Println(buf)
}

func printf(level Level, format string, args ...interface{}) {
	if level > stringToLevel(l.config.Level) {
		return
	}
	l.logger.SetPrefix(getPrefixByLevel(level))

	header := header()
	body := fmt.Sprintf(format, args...)

	buf := contentToBuffer(header, body)
	l.logger.Printf(buf.String())
}

func newLogger(conf LogConfig) *syslog.Logger {
	path := path.Join(conf.Dir, conf.FileName)
	fd, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		syslog.Fatal(err)
		os.Exit(1)
	}

	flag := syslog.Ldate | syslog.Ltime | syslog.Lshortfile | syslog.Lmicroseconds

	logger := syslog.New(fd, "", flag)

	return logger
}

func rotateDaemon() {
	for {
		time.Sleep(time.Second * 1)
		// 切割
		currentTime := getCurrentTime(l.config)
		if currentTime != l.rotateTime {
			os.Rename(l.config.FileName, l.config.FileName+fmt.Sprintf(".%s", l.rotateTime))
			// 重新设置logger与rotateTime
			// TODO：执行上一步后，下一步没执行前，其它goroutine会有问题吗？
			SetLogger(newLogger(l.config))
			SetRotateTime(currentTime)
		}

		// 删除
		files, err := ioutil.ReadDir(l.config.Dir)
		// 保留n天
		minKeepTime := time.Now().AddDate(0, 0, -l.config.KeepDays)
		reg := regexp.MustCompile("\\.log\\.20[0-9]{8}")
		if err == nil {
			for _, file := range files {
				if reg.FindString(file.Name()) != "" && shouldDel(file.Name(), minKeepTime) {
					os.Remove(path.Join(l.config.Dir, file.Name()))
				}
			}
		}
	}
}
