# Go Libary
<!-- vim-markdown-toc GFM -->

* [1. Log](#1-log)
    * [1.1 配置](#11-配置)
    * [1.2 使用](#12-使用)
* [2. Pool](#2-pool)
    * [2.1 配置](#21-配置)
    * [2.2 使用](#22-使用)
* [3. Client](#3-client)
    * [3.1 Http](#31-http)
        * [3.1.1 配置](#311-配置)
        * [3.1.2 使用](#312-使用)
    * [3.2 Redis](#32-redis)
        * [3.2.1 配置](#321-配置)
        * [3.2.2 使用](#322-使用)
    * [3.3 MySQL](#33-mysql)
        * [3.3.1 配置](#331-配置)
        * [3.3.2 使用](#332-使用)
    * [3.4 Zookeeper](#34-zookeeper)
        * [3.4.1 配置](#341-配置)
        * [3.4.2 使用](#342-使用)

<!-- vim-markdown-toc -->

### 1. Log

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

### 2. Pool

#### 2.1 配置

```
MaxIdle
MaxActive
IdleTimeoutS
Dial
TestOnBorrow
Wait
```

#### 2.2 使用
```
pool := New(
    MaxIdle,
    MaxActive,
    IdleTimeoutS,
    func() (Conn, error) {
        index := rand_gen.Intn(len(Servers))
        c, err := net.DialTimeout(
            "tcp",
            Servers[index],
            time.Duration(ConnTimeoutMs)*time.Millisecond,
        )
        return c, err
    },
    nil,
    true,
)
```

## 3. Client

### 3.1 Http

#### 3.1.1 配置
```
TimeoutMs           int // 总超时时间，单位毫秒
ConnectTimeoutMs    int // 连接超时时间，单位毫秒
KeepAlive           int // 长连接过期时间，单位秒
MaxIdleConnsPerHost int
```

#### 3.1.2 使用
```
client := client.New(TimeoutMs, ConnectTimeoutMs, KeepAlive, MaxIdleConnsPerhost)
resp, err = client.Get("http://www.baidu.com/s", map[string]interface{}{"wd": "beijing"})
resp, err = client.Post("http://www.baidu.com/s", map[string]interface{}{"wd": "beijing"})
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
client := &redis.Client{conf}
client.Init()
client.Set("hello", []byte("world"))
client.Get("hello")
```

### 3.3 MySQL

#### 3.3.1 配置

```
"Server": "127.0.0.1:3306",
"User": "root",
"Password": "",
"DataBase": "test",
"ConnTimeoutMs": 200,
"WriteTimeoutMs": 200,
"ReadTimeoutMs": 1000,
"MaxIdleConn": 50,
"MaxOpenConn": 200
```

#### 3.3.2 使用

```
client := &mysql.Client{conf}
client.Init()
client.DB.Table("users").First(&user)
```


### 3.4 Zookeeper

#### 3.4.1 配置

```
Servers          []string `json:"servers"`
ConnectTimeoutMs int      `json:"connect_timeout"`
```

#### 3.4.2 使用

```
client := &zookeeper.Client{conf}
client.Init()
content, stat, ch, err := client.GetW("/test")
```
