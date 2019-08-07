package mysql

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type User struct {
	Id        int64
	Name      string
	Age       int
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func TestSql(t *testing.T) {
	clients, err := Init("../../conf/mysql.conf")
	if err != nil {
		t.Fatal(err)
		os.Exit(1)
	}
	testDb := clients["test"].DB

	user := User{}
	testDb.Table("users").First(&user)
	fmt.Println(map[string]interface{}{
		"name":       user.Name,
		"created_at": user.CreatedAt,
	})
}
