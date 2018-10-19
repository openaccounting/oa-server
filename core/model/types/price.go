package types

import (
	"time"
)

type Price struct {
	Id       string    `json:"id"`
	OrgId    string    `json:"orgId"`
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
	Inserted time.Time `json:"inserted"`
	Updated  time.Time `json:"updated"`
	Price    float64   `json:"price"`
}
