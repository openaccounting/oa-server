package types

import (
	"time"
)

type Transaction struct {
	Id          string    `json:"id"`
	OrgId       string    `json:"orgId"`
	UserId      string    `json:"userId"`
	Date        time.Time `json:"date"`
	Inserted    time.Time `json:"inserted"`
	Updated     time.Time `json:"updated"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	Deleted     bool      `json:"deleted"`
	Splits      []*Split  `json:"splits"`
}

type Split struct {
	TransactionId string `json:"-"`
	AccountId     string `json:"accountId"`
	Amount        int64  `json:"amount"`
	NativeAmount  int64  `json:"nativeAmount"`
}
