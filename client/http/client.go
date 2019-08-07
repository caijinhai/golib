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
	TimeoutMs           int // 总超时时间，单位毫秒
	ConnectTimeoutMs    int // 连接超时时间，单位毫秒
	KeepAlive           int // 长连接过期时间，单位秒
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
	defaultClient.SetTimeout(DEAULT_TIMEOUTMS)
	defaultClient.SetConnectTimeout(DEAULT_CONNECT_TIMEOUTMS)
	defaultClient.SetKeepAlive(DEAULT_KEEPALIVES)
	defaultClient.SetMaxIdleConnsPerHost(10)
	defaultClient.SetHandler(defaultClient.NewHandler())
}

// 局部对象
func New(TimeoutMs int, ConnectTimeoutMs int, KeepAlive int, MaxIdleConnsPerhost int) *Client {
	client := &Client{}
	client.SetTimeout(client.TimeoutMs)
	client.SetConnectTimeout(client.ConnectTimeoutMs)
	client.SetKeepAlive(client.KeepAlive)
	client.SetMaxIdleConnsPerHost(MaxIdleConnsPerhost)
	client.SetHandler(client.NewHandler())
	return client
}

func (client *Client) NewHandler() *http.Client {
	return &http.Client{
		Timeout: time.Duration(client.TimeoutMs) * time.Millisecond, //总的超时
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(client.ConnectTimeoutMs) * time.Millisecond, //链接超时
				KeepAlive: time.Duration(client.KeepAlive) * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: client.MaxIdleConnsPerHost, //每个host连接池的大小
		},
	}
}

func (client *Client) SetTimeout(Timeout int) {
	client.TimeoutMs = Timeout
}

func (client *Client) SetConnectTimeout(Timeout int) {
	client.ConnectTimeoutMs = Timeout
}

func (client *Client) SetKeepAlive(KeepAlive int) {
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
