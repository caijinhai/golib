package helper

import (
	"bytes"
	"encoding/json"
	"net/url"
)

func Map2Buffer(mm map[string]interface{}) *bytes.Buffer {
	data, err := json.Marshal(mm)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer([]byte(data))
}

func Map2UrlParams(mm map[string]interface{}) string {
	vals := url.Values{}
	for key, val := range mm {
		vals.Add(key, val.(string))
	}
	return vals.Encode()
}
