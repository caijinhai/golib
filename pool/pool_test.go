package pool

import (
	"math/rand"
	"net"
	"testing"
	"time"
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

func TestPool(t *testing.T) {

	rand_gen := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := Client{
		ConnTimeoutMs:  50,
		WriteTimeoutMs: 200,
		ReadTimeoutMs:  200,
		IdleTimeoutS:   60,
		MaxIdle:        100,
		MaxActive:      200,
		Addrs:          []string{"127.0.0.1:6379"},
	}

	client.pool = NewPool(
		client.MaxIdle,
		client.MaxActive,
		func() (Conn, error) {
			index := rand_gen.Intn(len(client.Addrs))
			c, err := net.Dial("tcp", client.Addrs[index])
			if err == nil {
				err = c.SetWriteDeadline(time.Now().Add(time.Duration(client.WriteTimeoutMs) * time.Nanosecond))
				err = c.SetReadDeadline(time.Now().Add(time.Duration(client.ReadTimeoutMs) * time.Nanosecond))
			}
			return c, err
		},
	)

	conn, err := client.pool.Get()
	if err != nil {
		t.Fatal(err)
	}

	client.pool.Close(conn)
}
