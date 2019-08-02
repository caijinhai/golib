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

## 连接池

### 核心配置

```
MaxIdle
MaxActive
Dial
Wait
```

### 使用
```
import (
    "net"
    "time"
    "math/rand"
)
rand_gen := rand.New(rand.NewSource(time.Now().UnixNano()))

client.pool = pool.New(
    client.MaxIdle,
    client.MaxActive,
    client.IdleTimeoutS,
    func() (pool.Conn, error) {
        index := rand_gen.Intn(len(client.Servers))
        c, err := net.Dial("tcp", client.Servers[index])
        if err == nil {
            err = c.SetWriteDeadline(time.Now().Add(time.Duration(client.WriteTimeoutMs) * time.Nanosecond))
            err = c.SetReadDeadline(time.Now().Add(time.Duration(client.ReadTimeoutMs) * time.Nanosecond))
        }
        return c, err
    },
    true,
)
```
