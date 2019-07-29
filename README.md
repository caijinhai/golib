# Go Libary

```
├── client
│   ├── http
│   ├── mysql
│   ├── redis
│   ├── thrift
│   └── zookeeper
├── helper
│   ├── arr.go
│   ├── helper.go
│   ├── linklist.go
│   ├── str.go
│   └── time.go
├── log
│   └── log.go
└── utils
    ├── hashid.go
    ├── validator.go
    └── view.go
```

## 日志

### 配置文件
```
{
    "Type": "file",
    "Level": "DEBUG",
    "Dir": "./",
    "FileName": "test.log",
    "RotateByHour": true, // 按照小时分割
    "KeepDays": 7 // 保留7天
}
```

### 使用
```
log.Init("./log.conf") 
log.SetTraceId("请求id")
log.Debug(map[string]interface{}{
	"action": "test",
	"result": "success",
})
log.Debugf("xxxxx")
```