package http

import (
	"github.com/caijinlin/golib/helper"
	"github.com/caijinlin/golib/log"
	"net"
	"net/http"
	"time"
)

// 参考 https://github.com/nahid/gohttp

// 全局变量/局部变量

type Client struct {
	Timeout             time.Duration // 总超时时间
	ConnectTimeout      time.Duration // 连接超时时间
	KeepAlive           time.Duration // 长连接过期时间
	MaxIdleConnsPerHost int
	Handler             *http.Client
}

const (
	DEAULT_TIMEOUTMS         = 2000 // 2000 ms
	DEAULT_CONNECT_TIMEOUTMS = 200  // 200 ms
	DEAULT_KEEPALIVES        = 30   // 30 s
)

// 全局对象
var defaultClient = &Client{}

func init() {
	defaultClient.SetTimeout(time.Duration(DEAULT_TIMEOUTMS) * time.Millisecond)
	defaultClient.SetConnectTimeout(time.Duration(DEAULT_CONNECT_TIMEOUTMS) * time.Millisecond)
	defaultClient.SetKeepAlive(time.Duration(DEAULT_KEEPALIVES) * time.Second)
	defaultClient.SetMaxIdleConnsPerHost(10)
	defaultClient.SetHandler(defaultClient.NewHandler())
}

// 局部对象
func New(Timeout int, ConnectTimeout int, KeepAlive int, MaxIdleConnsPerhost int) *Client {
	client := &Client{}
	client.SetTimeout(time.Duration(Timeout) * time.Millisecond)
	client.SetConnectTimeout(time.Duration(ConnectTimeout) * time.Millisecond)
	client.SetKeepAlive(time.Duration(KeepAlive) * time.Second)
	client.SetMaxIdleConnsPerHost(MaxIdleConnsPerhost)
	client.SetHandler(client.NewHandler())
	return client
}

func (client *Client) NewHandler() *http.Client {
	return &http.Client{
		Timeout: client.Timeout, //总的超时
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   client.ConnectTimeout, //链接超时
				KeepAlive: client.KeepAlive,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: client.MaxIdleConnsPerHost, //每个host连接池的大小
		},
	}
}

func (client *Client) SetTimeout(Timeout time.Duration) {
	client.Timeout = Timeout
}

func (client *Client) SetConnectTimeout(Timeout time.Duration) {
	client.ConnectTimeout = Timeout
}

func (client *Client) SetKeepAlive(KeepAlive time.Duration) {
	client.KeepAlive = KeepAlive
}

func (client *Client) SetMaxIdleConnsPerHost(MaxIdle int) {
	client.MaxIdleConnsPerHost = MaxIdle
}

func (client *Client) SetHandler(Handler *http.Client) {
	client.Handler = Handler
}

// client API

func Get(url string, params map[string]interface{}) (*Response, error) {
	return defaultClient.do(http.MethodGet, url, params)
}

func Post(url string, params map[string]interface{}) (*Response, error) {
	return defaultClient.do(http.MethodPost, url, params)
}

func Put(url string, params map[string]interface{}) (*Response, error) {
	return defaultClient.do(http.MethodPut, url, params)
}

func Delete(url string, params map[string]interface{}) (*Response, error) {
	return defaultClient.do(http.MethodDelete, url, params)
}

func (client *Client) Get(url string, params map[string]interface{}) (*Response, error) {
	return client.do(http.MethodGet, url, params)
}

func (client *Client) Post(url string, params map[string]interface{}) (*Response, error) {
	return client.do(http.MethodPost, url, params)
}

func (client *Client) Put(url string, params map[string]interface{}) (*Response, error) {
	return client.do(http.MethodPut, url, params)
}

func (client *Client) Delete(url string, params map[string]interface{}) (*Response, error) {
	return client.do(http.MethodDelete, url, params)
}

/**
* 统一收敛调用入口，并加上耗时统计
 */
func (client *Client) do(method string, url string, params map[string]interface{}) (response *Response, err error) {

	start := time.Now()
	defer func() {
		cost := time.Now().Sub(start)
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		log.Info(map[string]interface{}{
			"action": "http_call",
			"url":    url,
			"method": method,
			"params": params,
			"cost":   helper.FormatDurationToMs(cost),
			"errmsg": errmsg,
		})
	}()

	if method == http.MethodGet {
		queryVals := helper.Map2UrlParams(params)
		if queryVals != "" {
			url += "?" + queryVals
		}
	}
	req := makeRequest(method, url, params)
	resp, err := client.Handler.Do(req)
	if err != nil {
		return nil, err
	}

	response = &Response{}
	response.resp = resp
	return response, err
}
