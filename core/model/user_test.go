package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/mocks"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TdUser struct {
	db.Datastore
	testNum int
}

func (td *TdUser) InsertUser(user *types.User) error {
	return nil
}

func (td *TdUser) UpdateUser(user *types.User) error {
	return nil
}

func TestCreateUser(t *testing.T) {

	// Id              string    `json:"id"`
	// Inserted        time.Time `json:"inserted"`
	// Updated         time.Time `json:"updated"`
	// FirstName       string    `json:"firstName"`
	// LastName        string    `json:"lastName"`
	// Email           string    `json:"email"`
	// Password        string    `json:"password"`
	// PasswordHash    string    `json:"-"`
	// AgreeToTerms    bool      `json:"agreeToTerms"`
	// PasswordReset   string    `json:"-"`
	// EmailVerified   bool      `json:"emailVerified"`
	// EmailVerifyCode string    `json:"-"`

	user := types.User{
		"0",
		time.Unix(0, 0),
		time.Unix(0, 0),
		"John",
		"Doe",
		"johndoe@email.com",
		"password",
		"",
		true,
		"",
		false,
		"",
	}

	badUser := types.User{
		"0",
		time.Unix(0, 0),
		time.Unix(0, 0),
		"John",
		"Doe",
		"",
		"password",
		"",
		true,
		"",
		false,
		"",
	}

	tests := map[string]struct {
		err  error
		user types.User
	}{
		"successful": {
			err:  nil,
			user: user,
		},
		"with error": {
			err:  errors.New("email required"),
			user: badUser,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		user := test.user

		mockBcrypt := new(mocks.Bcrypt)

		mockBcrypt.On("GetDefaultCost").Return(10)

		mockBcrypt.On("GenerateFromPassword", []byte(user.Password), 10).
			Return(make([]byte, 0), nil)

		model := NewModel(&TdUser{}, mockBcrypt, types.Config{})

		err := model.CreateUser(&user)

		assert.Equal(t, err, test.err)

		if err == nil {
			mockBcrypt.AssertExpectations(t)
		}
	}
}

func TestUpdateUser(t *testing.T) {

	user := types.User{
		"0",
		time.Unix(0, 0),
		time.Unix(0, 0),
		"John2",
		"Doe",
		"johndoe@email.com",
		"password",
		"",
		true,
		"",
		false,
		"",
	}

	badUser := types.User{
		"0",
		time.Unix(0, 0),
		time.Unix(0, 0),
		"John2",
		"Doe",
		"johndoe@email.com",
		"",
		"",
		true,
		"",
		false,
		"",
	}

	tests := map[string]struct {
		err  error
		user types.User
	}{
		"successful": {
			err:  nil,
			user: user,
		},
		"with error": {
			err:  errors.New("password required"),
			user: badUser,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		user := test.user

		mockBcrypt := new(mocks.Bcrypt)

		mockBcrypt.On("GetDefaultCost").Return(10)

		mockBcrypt.On("GenerateFromPassword", []byte(user.Password), 10).
			Return(make([]byte, 0), nil)

		model := NewModel(&TdUser{}, mockBcrypt, types.Config{})

		err := model.UpdateUser(&user)

		assert.Equal(t, err, test.err)

		if err == nil {
			mockBcrypt.AssertExpectations(t)
		}
	}
}
