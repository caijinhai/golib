package redis

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {

	redisClient, err := Init("../../conf/redis.conf")
	if err != nil {
		t.Fatal(err)
	}
	value, err := redisClient["order"].Get("test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(map[string]interface{}{
		"action": "getRedis",
		"value":  string(value),
		"err":    err,
	})
}
