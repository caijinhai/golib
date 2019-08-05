package http

import (
	"github.com/caijinlin/golib/helper"
	"net/http"
)

func makeRequest(method string, url string, params map[string]interface{}) *http.Request {

	request, err := http.NewRequest(method, url, helper.Map2Buffer(params))

	if err == nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("Connection", "keep-alive")
	}

	return request
}
