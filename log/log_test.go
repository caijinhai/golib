package log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {

	err := Init("../conf/log.conf")
	if err != nil {
		t.Fatal(err)
	}

	SetTraceId("afcc445911111")

	for {
		Debug(map[string]interface{}{
			"action": "test",
			"result": "success",
		})
		time.Sleep(time.Millisecond * 2)
		Infof("xxxxx")
		break
	}

}
