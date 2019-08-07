package mysql

import (
	"fmt"
)

func (client *Client) genDsn() string {

	paramstr := "charset=utf8&parseTime=True&loc=Local"
	if client.ConnTimeoutMs > 0 {
		paramstr += fmt.Sprintf("&timeout=%dms", client.ConnTimeoutMs)
	}
	if client.WriteTimeoutMs > 0 {
		paramstr += fmt.Sprintf("&writeTimeout=%dms", client.WriteTimeoutMs)
	}
	if client.ReadTimeoutMs > 0 {
		paramstr += fmt.Sprintf("&readTimeout=%dms", client.ReadTimeoutMs)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		client.User,
		client.Password,
		client.Server, // 127.0.0.1:3306
		client.DataBase,
		paramstr,
	)

	return dsn
}
