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
	err = redisClient["order"].Set("test_key1", []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	value, err := redisClient["order"].Get("test_key1")

	fmt.Println(map[string]interface{}{
		"action": "getRedis",
		"value":  string(value),
		"err":    err,
	})
}

func TestSentinel(t *testing.T) {

	redisClient, err := Init("../../conf/redis.conf")
	if err != nil {
		t.Fatal(err)
	}
	err = redisClient["rider"].Set("test_key2", []byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	value, err := redisClient["rider"].Get("test_key2")

	fmt.Println(map[string]interface{}{
		"action": "getRedis",
		"value":  string(value),
		"err":    err,
	})
}
