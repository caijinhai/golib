package http

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {

	resp, err := Get("http://www.baidu.com/s", map[string]interface{}{"wd": "beijing"})
	if err != nil {
		t.Fatal(err)
	}
	// body, _ := resp.GetBodyAsString()
	fmt.Println(resp.Protocol(), resp.GetStatusCode())
}

func TestPost(t *testing.T) {

	client := New(1000, 100, 30, 100)
	resp, err := client.Post("http://wwww.baidu.com", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Protocol(), resp.GetStatusCode())
}
