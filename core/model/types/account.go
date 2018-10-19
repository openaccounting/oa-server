package types

import (
	"time"
)

type Account struct {
	Id            string    `json:"id"`
	OrgId         string    `json:"orgId"`
	Inserted      time.Time `json:"inserted"`
	Updated       time.Time `json:"updated"`
	Name          string    `json:"name"`
	Parent        string    `json:"parent"`
	Currency      string    `json:"currency"`
	Precision     int       `json:"precision"`
	DebitBalance  bool      `json:"debitBalance"`
	Balance       *int64    `json:"balance"`
	NativeBalance *int64    `json:"nativeBalance"`
	ReadOnly      bool      `json:"readOnly"`
	HasChildren   bool      `json:"-"`
}

type AccountNode struct {
	Account  *Account
	Parent   *AccountNode
	Children []*AccountNode
}

func NewAccount() *Account {
	return &Account{Precision: 2}
}
