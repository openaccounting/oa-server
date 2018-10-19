package types

import (
	"time"
)

type User struct {
	Id              string    `json:"id"`
	Inserted        time.Time `json:"inserted"`
	Updated         time.Time `json:"updated"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	PasswordHash    string    `json:"-"`
	AgreeToTerms    bool      `json:"agreeToTerms"`
	PasswordReset   string    `json:"-"`
	EmailVerified   bool      `json:"emailVerified"`
	EmailVerifyCode string    `json:"-"`
}
