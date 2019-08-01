package pool

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/caijinlin/golib/log"
)

/**
* 各个client
 */
type Client struct {
	ConnTimeoutMs  int
	WriteTimeoutMs int
	ReadTimeoutMs  int
	IdleTimeoutS   int // 单位秒
	MaxIdle        int
	MaxActive      int
	Addrs          []string
	pool           *ConnPool
}

func init() {
	log.Init("../log/log.conf")
}

func TestPool(t *testing.T) {

	rand_gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := Client{
		ConnTimeoutMs:  50,
		WriteTimeoutMs: 200,
		ReadTimeoutMs:  200,
		IdleTimeoutS:   60,
		MaxIdle:        100,
		MaxActive:      1,
		Addrs:          []string{"39.156.66.14:80"},
	}

	client.pool = New(
		client.MaxIdle,
		client.MaxActive,
		client.IdleTimeoutS,
		func() (Conn, error) {
			index := rand_gen.Intn(len(client.Addrs))
			c, err := net.DialTimeout(
				"tcp",
				client.Addrs[index],
				time.Duration(client.ConnTimeoutMs)*time.Millisecond,
			)
			return c, err
		},
		true,
	)

	fmt.Println("准备获取连接conn1")
	conn1, err := client.pool.Get()
	fmt.Println("完成获取连接conn1")

	go func() {
		fmt.Println("准备释放conn1")
		time.Sleep(time.Duration(1) * time.Second)
		client.pool.Release(conn1)
		fmt.Println("完成释放conn1")
	}()
	if err != nil {
		fmt.Println(111)
		t.Fatal(err)
	}
	fmt.Println("准备获取连接conn2")
	conn2, err := client.pool.Get()
	fmt.Println("完成获取连接conn2")
	if err != nil {
		fmt.Println(222)
		t.Fatal(err)
	}
	client.pool.Release(conn2)
	log.Info(map[string]interface{}{
		"xxx": "Nice, you are great",
	})
}
