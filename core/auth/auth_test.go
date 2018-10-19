package auth

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TdUser struct {
	db.Datastore
	testNum int
}

func (td *TdUser) GetVerifiedUserByEmail(email string) (*types.User, error) {
	switch td.testNum {
	case 1:
		return td.GetVerifiedUserByEmail_1(email)
	case 2:
		return td.GetVerifiedUserByEmail_2(email)
	}

	return nil, errors.New("test error")
}

func (td *TdUser) GetVerifiedUserByEmail_1(email string) (*types.User, error) {
	return &types.User{
		"1",
		time.Unix(0, 0),
		time.Unix(0, 0),
		"John",
		"Doe",
		"johndoe@email.com",
		"password",
		"$2a$10$KrtvADe7jwrmYIe3GXFbNupOQaPIvyOKeng5826g4VGOD47TpAisG",
		true,
		"",
		false,
		"",
	}, nil
}

func (td *TdUser) GetVerifiedUserByEmail_2(email string) (*types.User, error) {
	return nil, errors.New("sql error")
}

func TestAuthenticateUser(t *testing.T) {
	tests := map[string]struct {
		err        error
		email      string
		password   string
		saltedHash string
		testNum    int
	}{
		"successful": {
			err:        nil,
			email:      "johndoe@email.com",
			password:   "password",
			saltedHash: "$2a$10$KrtvADe7jwrmYIe3GXFbNupOQaPIvyOKeng5826g4VGOD47TpAisG",
			testNum:    1,
		},
		"non-existing user": {
			err:        errors.New("Invalid email or password"),
			email:      "nouser@email.com",
			password:   "password",
			saltedHash: "",
			testNum:    2,
		},
		"wrong password": {
			err:        errors.New("Invalid email or password"),
			email:      "johndoe@email.com",
			password:   "bad",
			saltedHash: "$2a$10$KrtvADe7jwrmYIe3GXFbNupOQaPIvyOKeng5826g4VGOD47TpAisG",
			testNum:    1,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		authService := NewAuthService(&TdUser{testNum: test.testNum}, new(util.StandardBcrypt))

		_, err := authService.AuthenticateUser(test.email, test.password)

		assert.Equal(t, err, test.err)
	}
}
