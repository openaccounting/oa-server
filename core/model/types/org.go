package types

import (
	"time"
)

type Org struct {
	Id        string    `json:"id"`
	Inserted  time.Time `json:"inserted"`
	Updated   time.Time `json:"updated"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	Precision int       `json:"precision"`
	Timezone  string    `json:"timezone"`
}
