package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caijinlin/golib/client/mysql/driver"
	"github.com/caijinlin/golib/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"os"
)

type Client struct {
	ConnTimeoutMs  int `json:"ConnTimeoutMs"`
	WriteTimeoutMs int `json:"WriteTimeoutMs"`
	ReadTimeoutMs  int `json:"ReadTimeoutMs"`
	MaxIdle        int // 连接池中的最大连接数
	MaxActive      int // 最大活跃数

	Server   string `json:"Server"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DataBase string `json:"DataBase"`

	// 系统默认*sql.DB
	DB *gorm.DB // 应该是个interface
}

/**
* 通过配置文件生成client
 */
func Init(confFile string) (clients map[string]*Client, err error) {
	if res, err := ioutil.ReadFile(confFile); err != nil {
		err = errors.New("error opening conf file=" + confFile)
	} else {
		if err := json.Unmarshal(res, &clients); err != nil {
			msg := fmt.Sprintf("error parsing conf file=%s, err=%s", confFile, err.Error())
			err = errors.New(msg)
		}
	}

	if err != nil {
		return
	}
	for key, _ := range clients {
		clients[key].Init()
	}

	return
}

/**
* 使用gorm初始化client
**/
func (client *Client) Init() {
	var err error
	client.DB, err = driver.NewGorm(client.genDsn())
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
	// gorm->DB() = sql.DB
	client.DB.DB().SetMaxIdleConns(client.MaxIdle)
	client.DB.DB().SetMaxOpenConns(client.MaxActive)
}
