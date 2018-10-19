package types

import (
	"time"
)

type Invite struct {
	Id       string    `json:"id"`
	OrgId    string    `json:"orgId"`
	Inserted time.Time `json:"inserted"`
	Updated  time.Time `json:"updated"`
	Email    string    `json:"email"`
	Accepted bool      `json:"accepted"`
}
