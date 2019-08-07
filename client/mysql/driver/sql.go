package driver

import (
	"database/sql"
)

func NewSql(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	return db, err
}
