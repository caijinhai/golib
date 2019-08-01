package redis

import (
	"github.com/caijinlin/golib/log"
	"github.com/caijinlin/golib/pool"
	redislib "github.com/garyburd/redigo/redis"
	"math/rand"
	"net"
	"time"
)

type Client struct {
	ConnTimeoutMs  int // 单位毫秒
	WriteTimeoutMs int // 单位毫秒
	ReadTimeoutMs  int // 单位毫秒
	IdleTimeoutS   int // 单位秒
	MaxIdle        int // 连接池中的最大连接数
	MaxActive      int // 最大活跃数
	Addrs          []string
	Password       string
	Db             int
	pool           *pool.ConnPool
}

/**
* 通过配置文件转化为client，然后init，方便调用者
**/
func (client *Client) Init() {
	client.initPool()
}

func (client *Client) Get(key string) (value []byte, err error) {
	value, err = client.Do("GET", key)
	return
}

func (client *Client) Set(key string, value []byte) (err error) {
	_, err = client.Do("SET", key, value)
	return
}

func (client *Client) Lock(key string, expire_ms int) (token string, err error) {
	token = string(rand.Int())
	_, err = client.Do("SET", key, []byte(token), "PX", expire_ms, "NX")

	return
}

func (client *Client) Unlock(key string, value string) (err error) {
	var script = redislib.NewScript(2,
		`if redis.call("GET", KEYS[1]) == KEYS[2] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end`)
	_, err = client.DoScript(script, key, value)

	return
}

func (client *Client) DoScript(scirpt *redislib.Script, args ...interface{}) (reply []byte, err error) {
	conn, err := client.pool.Get()
	if err != nil {
		log.Warning(map[string]interface{}{
			"action": "poolGetConn",
			"err":    err,
		})
	}
	defer client.pool.Release(conn)
	redisConn, _ := conn.(redislib.Conn)
	reply, err = redislib.Bytes(scirpt.Do(redisConn, args...))
	return
}

func (client *Client) Do(commandName string, args ...interface{}) (reply []byte, err error) {
	conn, err := client.pool.Get()
	if err != nil {
		log.Warning(map[string]interface{}{
			"action": "poolGetConn",
			"err":    err,
		})
		return
	}
	defer client.pool.Release(conn)
	redisConn, _ := conn.(redislib.Conn)
	reply, err = redislib.Bytes(redisConn.Do(commandName, args...))
	return
}

func (client *Client) initPool() {

	rand_gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	client.pool = pool.New(
		client.MaxIdle,
		client.MaxActive,
		client.IdleTimeoutS,
		func() (pool.Conn, error) {
			index := rand_gen.Intn(len(client.Addrs))
			// 网络连接
			netConn, err := net.DialTimeout(
				"tcp",
				client.Addrs[index],
				time.Duration(client.ConnTimeoutMs)*time.Millisecond,
			)
			// redis连接
			redisConn := redislib.NewConn(
				netConn,
				time.Duration(client.ReadTimeoutMs)*time.Millisecond,
				time.Duration(client.WriteTimeoutMs)*time.Millisecond,
			)
			if client.Password != "" {
				if _, err := redisConn.Do("AUTH", client.Password); err != nil {
					netConn.Close()
					return nil, err
				}
			}
			if client.Db > 0 {
				if _, err := redisConn.Do("SELECT", client.Db); err != nil {
					netConn.Close()
					return nil, err
				}
			}
			return redisConn, err
		},
		true,
	)
}
