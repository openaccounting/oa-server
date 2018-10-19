package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type TdAccount struct {
	db.Datastore
	mock.Mock
}

func (td *TdAccount) GetPermissionedAccountIds(userId string, orgId string, tokenId string) ([]string, error) {
	// User has permission to only "Assets" account
	return []string{"2"}, nil
}

func (td *TdAccount) GetAccountsByOrgId(orgId string) ([]*types.Account, error) {
	args := td.Called(orgId)
	return args.Get(0).([]*types.Account), args.Error(1)
}

func (td *TdAccount) InsertAccount(account *types.Account) error {
	return nil
}

func (td *TdAccount) UpdateAccount(account *types.Account) error {
	return nil
}

func (td *TdAccount) AddBalance(account *types.Account, date time.Time) error {
	return nil
}

func (td *TdAccount) AddNativeBalanceNearestInTime(account *types.Account, date time.Time) error {
	return nil
}

func (td *TdAccount) AddNativeBalanceCost(account *types.Account, date time.Time) error {
	return nil
}

func (td *TdAccount) AddBalances(accounts []*types.Account, date time.Time) error {
	balance := int64(1000)
	for _, account := range accounts {
		account.Balance = &balance
	}

	return nil
}

func (td *TdAccount) AddNativeBalancesNearestInTime(accounts []*types.Account, date time.Time) error {
	balance := int64(1000)
	for _, account := range accounts {
		account.NativeBalance = &balance
	}

	return nil
}

func (td *TdAccount) AddNativeBalancesCost(accounts []*types.Account, date time.Time) error {
	balance := int64(1000)
	for _, account := range accounts {
		account.NativeBalance = &balance
	}

	return nil
}

func (td *TdAccount) GetSplitCountByAccountId(id string) (int64, error) {
	args := td.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (td *TdAccount) GetChildCountByAccountId(id string) (int64, error) {
	args := td.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

func (td *TdAccount) DeleteAccount(id string) error {
	return nil
}

func (td *TdAccount) GetOrgUserIds(id string) ([]string, error) {
	return []string{"1"}, nil
}

func (td *TdAccount) GetAccount(id string) (*types.Account, error) {
	return &types.Account{}, nil
}

func getTestAccounts() []*types.Account {
	return []*types.Account{
		&types.Account{
			Id:           "2",
			OrgId:        "1",
			Name:         "Assets",
			Parent:       "1",
			Currency:     "USD",
			Precision:    2,
			DebitBalance: true,
		},
		&types.Account{
			Id:           "3",
			OrgId:        "1",
			Name:         "Current Assets",
			Parent:       "2",
			Currency:     "USD",
			Precision:    2,
			DebitBalance: true,
		},
		&types.Account{
			Id:           "1",
			OrgId:        "1",
			Name:         "Root",
			Parent:       "",
			Currency:     "USD",
			Precision:    2,
			DebitBalance: true,
		},
	}
}

func TestCreateAccount(t *testing.T) {
	tests := map[string]struct {
		err     error
		account *types.Account
	}{
		"success": {
			err: nil,
			account: &types.Account{
				Id:           "1",
				OrgId:        "1",
				Name:         "Cash",
				Parent:       "3",
				Currency:     "USD",
				Precision:    2,
				DebitBalance: true,
			},
		},
		"permission error": {
			err: errors.New("user does not have permission to access account 1"),
			account: &types.Account{
				Id:           "1",
				OrgId:        "1",
				Name:         "Cash",
				Parent:       "1",
				Currency:     "USD",
				Precision:    2,
				DebitBalance: true,
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdAccount{}
		td.On("GetAccountsByOrgId", "1").Return(getTestAccounts(), nil)

		model := NewModel(td, nil, types.Config{})

		err := model.CreateAccount(test.account, "1")
		assert.Equal(t, test.err, err)
	}
}

func TestUpdateAccount(t *testing.T) {
	tests := map[string]struct {
		err     error
		account *types.Account
	}{
		"success": {
			err: nil,
			account: &types.Account{
				Id:           "3",
				OrgId:        "1",
				Name:         "Current Assets2",
				Parent:       "2",
				Currency:     "USD",
				Precision:    2,
				DebitBalance: true,
			},
		},
		"error": {
			err: errors.New("account cannot be its own parent"),
			account: &types.Account{
				Id:           "3",
				OrgId:        "1",
				Name:         "Current Assets",
				Parent:       "3",
				Currency:     "USD",
				Precision:    2,
				DebitBalance: true,
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdAccount{}
		td.On("GetAccountsByOrgId", "1").Return(getTestAccounts(), nil)

		model := NewModel(td, nil, types.Config{})

		err := model.UpdateAccount(test.account, "1")
		assert.Equal(t, test.err, err)

		if err == nil {
			td.AssertExpectations(t)
		}
	}
}

func TestDeleteAccount(t *testing.T) {
	tests := map[string]struct {
		err       error
		accountId string
		count     int64
	}{
		"success": {
			err:       nil,
			accountId: "3",
			count:     0,
		},
		"error": {
			err:       errors.New("Cannot delete an account that has transactions"),
			accountId: "3",
			count:     1,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdAccount{}
		td.On("GetAccountsByOrgId", "1").Return(getTestAccounts(), nil)
		td.On("GetSplitCountByAccountId", test.accountId).Return(test.count, nil)
		td.On("GetChildCountByAccountId", test.accountId).Return(test.count, nil)

		model := NewModel(td, nil, types.Config{})

		err := model.DeleteAccount(test.accountId, "1", "1")
		assert.Equal(t, test.err, err)

		if err == nil {
			td.AssertExpectations(t)
		}
	}
}

func TestGetAccounts(t *testing.T) {
	tests := map[string]struct {
		err error
	}{
		"success": {
			err: nil,
		},
		// "error": {
		// 	err: errors.New("db error"),
		// },
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdAccount{}
		td.On("GetAccountsByOrgId", "1").Return(getTestAccounts(), test.err)

		model := NewModel(td, nil, types.Config{})

		accounts, err := model.GetAccounts("1", "1", "")

		assert.Equal(t, test.err, err)

		if err == nil {
			td.AssertExpectations(t)
			assert.Equal(t, 3, len(accounts))
			assert.Equal(t, false, accounts[0].ReadOnly)
			assert.Equal(t, false, accounts[1].ReadOnly)
			assert.Equal(t, true, accounts[2].ReadOnly)
		}
	}
}

func TestGetAccountsWithBalances(t *testing.T) {
	tests := map[string]struct {
		err error
	}{
		"success": {
			err: nil,
		},
		"error": {
			err: errors.New("db error"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdAccount{}
		td.On("GetAccountsByOrgId", "1").Return(getTestAccounts(), test.err)

		model := NewModel(td, nil, types.Config{})

		accounts, err := model.GetAccountsWithBalances("1", "1", "", time.Now())

		assert.Equal(t, test.err, err)

		if err == nil {
			td.AssertExpectations(t)
			assert.Equal(t, 3, len(accounts))
			assert.Equal(t, false, accounts[0].ReadOnly)
			assert.Equal(t, false, accounts[1].ReadOnly)
			assert.Equal(t, true, accounts[2].ReadOnly)

			assert.Equal(t, int64(1000), *accounts[0].Balance)
			assert.Equal(t, int64(1000), *accounts[1].Balance)

			assert.Equal(t, int64(1000), *accounts[0].NativeBalance)
			assert.Equal(t, int64(1000), *accounts[1].NativeBalance)
		}
	}
}
