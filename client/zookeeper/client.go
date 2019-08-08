package zookeeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caijinlin/golib/log"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
	"time"
)

type Client struct {
	Servers          []string `json:"servers"`
	ConnectTimeoutMs int      `json:"connect_timeout"`
	ZkConn           *zk.Conn
}

func Init(confFile string) (client Client, err error) {
	if res, e := ioutil.ReadFile(confFile); e != nil {
		err = errors.New("error opening conf file=" + confFile)
	} else {
		if err = json.Unmarshal(res, &client); err != nil {
			err = errors.New(fmt.Sprintf("error parsing conf file=%s, err=%s", confFile, err.Error()))
		}
	}

	if err != nil {
		return
	}

	client.Init()
	return
}

func (client *Client) Init() {
	zkConn, _, err := zk.Connect(
		client.Servers,
		time.Duration(client.ConnectTimeoutMs)*time.Millisecond,
	)
	if err != nil {
		errmsg := fmt.Sprintf("Connecting zookeeper failed: %s", err.Error())
		log.Error(map[string]interface{}{"errmsg": errmsg})
	}
	client.ZkConn = zkConn
}

/**
* zk client API
 */

func (client *Client) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return client.ZkConn.Create(path, data, flags, acl)
}

func (client *Client) Delete(path string, version int32) error {
	return client.ZkConn.Delete(path, version)
}

func (client *Client) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return client.ZkConn.Set(path, data, version)
}

func (client *Client) Exists(path string, watch bool) (bool, *zk.Stat, error) {
	return client.ZkConn.Exists(path)
}

func (client *Client) ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error) {
	return client.ZkConn.ExistsW(path)
}

func (client *Client) Get(path string) ([]byte, *zk.Stat, error) {
	return client.ZkConn.Get(path)
}

func (client *Client) GetW(path string) ([]byte, *zk.Stat, <-chan zk.Event, error) {
	return client.ZkConn.GetW(path)
}

func (client *Client) Children(path string) ([]string, *zk.Stat, error) {
	return client.ZkConn.Children(path)
}

func (client *Client) ChildrenW(path string) ([]string, *zk.Stat, <-chan zk.Event, error) {
	return client.ZkConn.ChildrenW(path)
}
