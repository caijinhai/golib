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

## 1. Log

#### 1.1 配置
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

#### 1.2 使用
```
log.Init("./log.conf") 
log.SetTraceId("请求id")
log.Debug(map[string]interface{}{
	"action": "test",
	"result": "success",
})
log.Debugf("xxxxx")
```

## 2. Pool

#### 2.1 配置

```
MaxIdle
MaxActive
Dial
TestOnBorrow
Wait
```

#### 2.2 使用
```
rand_gen := rand.New(rand.NewSource(time.Now().UnixNano()))

client.pool = pool.New(
        client.MaxIdle,
        client.MaxActive,
        client.IdleTimeoutS,
        func() (pool.Conn, error) {
            index := rand_gen.Intn(len(client.Servers))
            redisConn, err := client.DialConn(client.Servers[index])
            return redisConn, err
        },
        func(c pool.Conn) error {
            conn, _ := c.(redislib.Conn)
            _, err := conn.Do("PING")
            return err
        },
        true,
    )
```

## 3. Client

### 3.1 Http

#### 3.1.1 配置
```
Timeout             time.Duration // 总超时时间
ConnectTimeout      time.Duration // 连接超时时间
KeepAlive           time.Duration // 长连接过期时间
MaxIdleConnsPerHost int // 每个host池子连接数
```

#### 3.1.2 使用
```
resp, err := client.Get("http://www.baidu.com/s", map[string]interface{}{"wd": "beijing"})
resp, err := client.Post("http://www.baidu.com/s", map[string]interface{}{"wd": "beijing"})
status := resp.GetStatusCode()
body, _ := resp.GetBodyAsString()
```

### 3.2 Redis

#### 3.2.1 配置

```
"SentinelServers": ["127.0.0.1:26379", "127.0.0.1:26380", "127.0.0.1:26381"],
"RedisSet": "api",
"Db":0,
"Servers": ["127.0.0.1:6379", "127.0.0.1:6380"],
"ConnTimeoutMs": 300,
"WriteTimeoutMs": 300,
"ReadTimeoutMs": 300,
"MaxIdle": 100,
"MaxActive": 200,
"IdleTimeoutS": 60
```

#### 3.2.2 使用

```
client := &redis.Client{}
client.Init()
client.Set("hello", []byte("world"))
client.Get("hello")
```