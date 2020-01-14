package types

import (
	"time"
)

type Budget struct {
	OrgId    string        `json:"orgId"`
	Inserted time.Time     `json:"inserted"`
	Items    []*BudgetItem `json:"items"`
}

type BudgetItem struct {
	OrgId     string `json:"-"`
	AccountId string `json:"accountId"`
	Amount    int64  `json:"amount"`
}
