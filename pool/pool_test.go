package pool

import (
	"fmt"
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
		ConnTimeoutMs:  100,
		WriteTimeoutMs: 200,
		ReadTimeoutMs:  200,
		IdleTimeoutS:   60,
		MaxIdle:        100,
		MaxActive:      200,
		Addrs:          []string{"39.156.66.14:80"},
	}

	client.pool = New(
		client.MaxIdle,
		client.MaxActive,
		client.IdleTimeoutS,
		func() (Conn, error) {
			index := rand_gen.Intn(len(client.Addrs))
			dialer := net.Dialer{Timeout: time.Duration(client.ConnTimeoutMs) * time.Millisecond}
			c, err := dialer.Dial("tcp", client.Addrs[index])
			if err == nil {
				c.SetWriteDeadline(time.Now().Add(time.Duration(client.WriteTimeoutMs) * time.Millisecond))
				c.SetReadDeadline(time.Now().Add(time.Duration(client.ReadTimeoutMs) * time.Millisecond))
			}
			return c, err
		},
	)

	conn, err := client.pool.Get()
	fmt.Println(conn)
	if err != nil {
		t.Fatal(err)
	}

	client.pool.Release(conn)
}
