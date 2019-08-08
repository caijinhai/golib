package zookeeper

import (
	"fmt"
	"os"
	"testing"
)

func TestZkApi(t *testing.T) {
	client, err := Init("../../conf/zookeeper.conf")
	if err != nil {
		t.Fatal(err)
		os.Exit(1)
	}
	content, stat, ch, err := client.GetW("/test")
	e := <-ch
	fmt.Println(content, stat, e, err)
}
