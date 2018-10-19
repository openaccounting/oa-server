package model

import (
	"errors"
	"fmt"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/ws"
	"time"
)

type TransactionInterface interface {
	CreateTransaction(*types.Transaction) error
	UpdateTransaction(string, *types.Transaction) error
	GetTransactionsByAccount(string, string, string, *types.QueryOptions) ([]*types.Transaction, error)
	GetTransactionsByOrg(string, string, *types.QueryOptions) ([]*types.Transaction, error)
	DeleteTransaction(string, string, string) error
}

func (model *Model) CreateTransaction(transaction *types.Transaction) (err error) {
	err = model.checkSplits(transaction)

	if err != nil {
		return
	}

	if transaction.Id == "" {
		return errors.New("id required")
	}

	transaction.Inserted = time.Now()
	transaction.Updated = time.Now()

	if transaction.Date.IsZero() {
		transaction.Date = transaction.Inserted
	}

	err = model.db.InsertTransaction(transaction)

	if err != nil {
		return
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access transaction
	userIds, err2 := model.db.GetOrgUserIds(transaction.OrgId)

	if err2 == nil {
		ws.PushTransaction(transaction, userIds, "create")
	}

	return
}

func (model *Model) UpdateTransaction(oldId string, transaction *types.Transaction) (err error) {
	err = model.checkSplits(transaction)

	if err != nil {
		return
	}

	if oldId == "" || transaction.Id == "" {
		return errors.New("id required")
	}

	// Get original transaction
	original, err := model.getTransactionById(oldId)

	if err != nil {
		return
	}

	transaction.Updated = time.Now()
	transaction.Inserted = original.Inserted

	// We used to compare splits and if they hadn't changed just do an update
	// on the transaction. The problem is then the updated field gets out of sync
	// between the tranaction and its splits.
	// It needs to be in sync for getTransactionsByOrg() to work correctly with pagination

	// Delete old transaction and insert a new one
	transaction.Inserted = transaction.Updated
	err = model.db.DeleteAndInsertTransaction(oldId, transaction)

	if err != nil {
		return
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access transaction
	userIds, err2 := model.db.GetOrgUserIds(transaction.OrgId)

	if err2 == nil {
		ws.PushTransaction(original, userIds, "delete")
		ws.PushTransaction(transaction, userIds, "create")
	}

	return
}

func (model *Model) GetTransactionsByAccount(orgId string, userId string, accountId string, options *types.QueryOptions) ([]*types.Transaction, error) {
	userAccounts, err := model.GetAccounts(orgId, userId, "")

	if err != nil {
		return nil, err
	}

	if !model.accountsContainWriteAccess(userAccounts, accountId) {
		return nil, errors.New(fmt.Sprintf("%s %s", "user does not have permission to access account", accountId))
	}

	return model.db.GetTransactionsByAccount(accountId, options)

}

func (model *Model) GetTransactionsByOrg(orgId string, userId string, options *types.QueryOptions) ([]*types.Transaction, error) {
	userAccounts, err := model.GetAccounts(orgId, userId, "")

	if err != nil {
		return nil, err
	}

	var accountIds []string
	for _, account := range userAccounts {
		accountIds = append(accountIds, account.Id)
	}

	return model.db.GetTransactionsByOrg(orgId, options, accountIds)
}

func (model *Model) DeleteTransaction(id string, userId string, orgId string) (err error) {
	transaction, err := model.getTransactionById(id)

	if err != nil {
		return
	}

	userAccounts, err := model.GetAccounts(orgId, userId, "")

	if err != nil {
		return
	}

	for _, split := range transaction.Splits {
		if !model.accountsContainWriteAccess(userAccounts, split.AccountId) {
			return errors.New(fmt.Sprintf("%s %s", "user does not have permission to access account", split.AccountId))
		}
	}

	err = model.db.DeleteTransaction(id)

	if err != nil {
		return
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access transaction
	userIds, err2 := model.db.GetOrgUserIds(transaction.OrgId)

	if err2 == nil {
		ws.PushTransaction(transaction, userIds, "delete")
	}

	return
}

func (model *Model) getTransactionById(id string) (*types.Transaction, error) {
	// TODO if this is made public, make a separate version that checks permission
	return model.db.GetTransactionById(id)
}

func (model *Model) checkSplits(transaction *types.Transaction) (err error) {
	if len(transaction.Splits) < 2 {
		return errors.New("at least 2 splits are required")
	}

	org, err := model.GetOrg(transaction.OrgId, transaction.UserId)

	if err != nil {
		return
	}

	userAccounts, err := model.GetAccounts(transaction.OrgId, transaction.UserId, "")

	if err != nil {
		return
	}

	var amount int64 = 0

	for _, split := range transaction.Splits {
		if !model.accountsContainWriteAccess(userAccounts, split.AccountId) {
			return errors.New(fmt.Sprintf("%s %s", "user does not have permission to access account", split.AccountId))
		}

		account := model.getAccountFromList(userAccounts, split.AccountId)

		if account.HasChildren == true {
			return errors.New("Cannot use parent account for split")
		}

		if account.Currency == org.Currency && split.NativeAmount != split.Amount {
			return errors.New("nativeAmount must equal amount for native currency splits")
		}

		amount += split.NativeAmount
	}

	if amount != 0 {
		return errors.New("splits must add up to 0")
	}

	return
}
