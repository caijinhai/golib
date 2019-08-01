package redis

import (
	"fmt"
	"testing"
)

func init() {

}

func TestClient(t *testing.T) {

	// TODO 通过配置文件转化
	client := Client{
		ConnTimeoutMs:  50,
		WriteTimeoutMs: 200,
		ReadTimeoutMs:  200,
		IdleTimeoutS:   60,
		MaxIdle:        100,
		MaxActive:      1,
		Addrs:          []string{"127.0.0.1:6379"},
		Db:             0,
		Password:       "",
	}
	client.Init()

	value, err := client.Get("test")
	fmt.Println(map[string]interface{}{
		"action": "getRedis",
		"value":  string(value),
		"err":    err,
	})
}
