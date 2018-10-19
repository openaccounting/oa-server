package types

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

type ApiKey struct {
	Id       string         `json:"id"`
	Inserted time.Time      `json:"inserted"`
	Updated  time.Time      `json:"updated"`
	UserId   string         `json:"userId"`
	Label    string         `json:"label"`
	Deleted  mysql.NullTime `json:"-"` // Can we marshal this correctly?
}
