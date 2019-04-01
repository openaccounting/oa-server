package model

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/ws"
)

type AccountInterface interface {
	CreateAccount(account *types.Account, userId string) error
	UpdateAccount(account *types.Account, userId string) error
	DeleteAccount(id string, userId string, orgId string) error
	GetAccounts(orgId string, userId string, tokenId string) ([]*types.Account, error)
	GetAccountsWithBalances(orgId string, userId string, tokenId string, date time.Time) ([]*types.Account, error)
	GetAccount(orgId, accId, userId, tokenId string) (*types.Account, error)
	GetAccountWithBalance(orgId, accId, userId, tokenId string, date time.Time) (*types.Account, error)
}

type ByName []*types.Account

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (model *Model) CreateAccount(account *types.Account, userId string) (err error) {
	if account.Id == "" {
		return errors.New("id required")
	}

	if account.OrgId == "" {
		return errors.New("orgId required")
	}

	if account.Name == "" {
		return errors.New("name required")
	}

	if account.Currency == "" {
		return errors.New("currency required")
	}

	userAccounts, err := model.GetAccounts(account.OrgId, userId, "")

	if err != nil {
		return
	}

	if !model.accountsContainWriteAccess(userAccounts, account.Parent) {
		return errors.New(fmt.Sprintf("%s %s", "user does not have permission to access account", account.Parent))
	}

	err = model.db.InsertAccount(account)

	if err != nil {
		return
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access account
	userIds, err2 := model.db.GetOrgUserIds(account.OrgId)

	if err2 == nil {
		ws.PushAccount(account, userIds, "create")
	}

	return
}

func (model *Model) UpdateAccount(account *types.Account, userId string) (err error) {
	if account.Id == "" {
		return errors.New("id required")
	}

	if account.OrgId == "" {
		return errors.New("orgId required")
	}

	if account.Name == "" {
		return errors.New("name required")
	}

	if account.Currency == "" {
		return errors.New("currency required")
	}

	if account.Parent == account.Id {
		return errors.New("account cannot be its own parent")
	}

	userAccounts, err := model.GetAccounts(account.OrgId, userId, "")

	if err != nil {
		return
	}

	if !model.accountsContainWriteAccess(userAccounts, account.Parent) {
		return errors.New(fmt.Sprintf("%s %s", "user does not have permission to access account", account.Parent))
	}

	err = model.db.UpdateAccount(account)

	if err != nil {
		return
	}

	err = model.db.AddBalance(account, time.Now())

	if err != nil {
		return
	}

	err = model.db.AddNativeBalanceCost(account, time.Now())

	if err != nil {
		return
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access account
	userIds, err2 := model.db.GetOrgUserIds(account.OrgId)

	if err2 == nil {
		ws.PushAccount(account, userIds, "update")
	}

	return
}

func (model *Model) DeleteAccount(id string, userId string, orgId string) (err error) {
	// TODO make sure user is part of org

	// check to make sure user has permission
	userAccounts, err := model.GetAccounts(orgId, userId, "")

	if err != nil {
		return
	}

	if !model.accountsContainWriteAccess(userAccounts, id) {
		return errors.New(fmt.Sprintf("%s %s", "user does not have permission to access account", id))
	}

	// don't allow deleting of accounts that have transactions or child accounts
	count, err := model.db.GetSplitCountByAccountId(id)

	if err != nil {
		return
	}

	if count != 0 {
		return errors.New("Cannot delete an account that has transactions")
	}

	count, err = model.db.GetChildCountByAccountId(id)

	if err != nil {
		return
	}

	if count != 0 {
		return errors.New("Cannot delete an account that has children")
	}

	account, err := model.db.GetAccount(id)

	if err != nil {
		return
	}

	err = model.db.DeleteAccount(id)

	if err != nil {
		return
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access account
	userIds, err2 := model.db.GetOrgUserIds(account.OrgId)

	if err2 == nil {
		ws.PushAccount(account, userIds, "delete")
	}

	return
}

func (model *Model) getAccounts(orgId string, userId string, tokenId string, date time.Time, withBalances bool) ([]*types.Account, error) {
	permissionedAccounts, err := model.db.GetPermissionedAccountIds(orgId, userId, "")
	if err != nil {
		return nil, err
	}

	var allAccounts []*types.Account

	if withBalances == true {
		allAccounts, err = model.getAllAccountsWithBalances(orgId, date)
	} else {
		allAccounts, err = model.getAllAccounts(orgId)
	}

	if err != nil {
		return nil, err
	}

	accountMap := model.makeAccountMap(allAccounts)
	writeAccessMap := make(map[string]*types.Account)
	readAccessMap := make(map[string]*types.Account)

	for _, accountId := range permissionedAccounts {
		writeAccessMap[accountId] = accountMap[accountId].Account

		// parents are read only
		parents := model.getParents(accountId, accountMap)

		for _, parentAccount := range parents {
			readAccessMap[parentAccount.Id] = parentAccount
		}

		// top level accounts are initially read only unless user has permission
		topLevelAccounts := model.getTopLevelAccounts(accountMap)

		for _, topLevelAccount := range topLevelAccounts {
			readAccessMap[topLevelAccount.Id] = topLevelAccount
		}

		// Children have write access
		children := model.getChildren(accountId, accountMap)

		for _, childAccount := range children {
			writeAccessMap[childAccount.Id] = childAccount
		}
	}

	filtered := make([]*types.Account, 0)

	for _, account := range writeAccessMap {
		filtered = append(filtered, account)
	}

	for id, account := range readAccessMap {
		_, ok := writeAccessMap[id]

		if ok == false {
			account.ReadOnly = true
			filtered = append(filtered, account)
		}
	}

	// TODO sort by inserted
	sort.Sort(ByName(filtered))

	return filtered, nil
}

func (model *Model) getAccount(orgId, accId, userId, tokenId string, date time.Time, withBalances bool) (*types.Account, error) {
	accounts, err := model.getAccounts(orgId, userId, tokenId, date, withBalances)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.Id == accId {
			return account, nil
		}
	}

	return nil, nil
}

func (model *Model) GetAccounts(orgId string, userId string, tokenId string) ([]*types.Account, error) {
	return model.getAccounts(orgId, userId, tokenId, time.Time{}, false)
}

func (model *Model) GetAccountsWithBalances(orgId string, userId string, tokenId string, date time.Time) ([]*types.Account, error) {
	return model.getAccounts(orgId, userId, tokenId, date, true)
}

func (model *Model) GetAccount(orgId, accId, userId, tokenId string) (*types.Account, error) {
	return model.getAccount(orgId, accId, userId, tokenId, time.Time{}, false)
}

func (model *Model) GetAccountWithBalance(orgId, accId, userId, tokenId string, date time.Time) (*types.Account, error) {
	return model.getAccount(orgId, accId, userId, tokenId, date, true)
}

func (model *Model) getAllAccounts(orgId string) ([]*types.Account, error) {
	return model.db.GetAccountsByOrgId(orgId)
}

func (model *Model) getAllAccountsWithBalances(orgId string, date time.Time) ([]*types.Account, error) {
	accounts, err := model.db.GetAccountsByOrgId(orgId)

	if err != nil {
		return nil, err
	}

	err = model.db.AddBalances(accounts, date)

	if err != nil {
		return nil, err
	}

	err = model.db.AddNativeBalancesCost(accounts, date)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (model *Model) makeAccountMap(accounts []*types.Account) map[string]*types.AccountNode {
	m := make(map[string]*types.AccountNode)

	for _, account := range accounts {
		m[account.Id] = &types.AccountNode{
			Account:  account,
			Parent:   nil,
			Children: nil,
		}
	}

	for _, account := range accounts {
		m[account.Id].Parent = m[account.Parent]

		if value, ok := m[account.Parent]; ok {
			value.Children = append(value.Children, m[account.Id])
			value.Account.HasChildren = true
		}
	}

	return m
}

func (model *Model) getChildren(parentId string, accountMap map[string]*types.AccountNode) []*types.Account {
	if _, ok := accountMap[parentId]; !ok {
		return nil
	}

	children := make([]*types.Account, 0)

	for _, childAccountNode := range accountMap[parentId].Children {
		children = append(children, childAccountNode.Account)
		grandChildren := model.getChildren(childAccountNode.Account.Id, accountMap)
		children = append(children, grandChildren...)
	}

	return children
}

func (model *Model) getParents(accountId string, accountMap map[string]*types.AccountNode) []*types.Account {
	node, ok := accountMap[accountId]

	if !ok {
		return nil
	}

	if node.Parent == nil {
		return make([]*types.Account, 0)
	}

	parents := model.getParents(node.Parent.Account.Id, accountMap)
	return append(parents, node.Parent.Account)
}

func (model *Model) accountsContainWriteAccess(accounts []*types.Account, accountId string) bool {
	for _, account := range accounts {
		if account.Id == accountId && !account.ReadOnly {
			return true
		}
	}
	return false
}

func (model *Model) getAccountFromList(accounts []*types.Account, accountId string) *types.Account {
	for _, account := range accounts {
		if account.Id == accountId {
			return account
		}
	}
	return nil
}

func (model *Model) getTopLevelAccounts(accountMap map[string]*types.AccountNode) []*types.Account {
	accounts := make([]*types.Account, 0)

	for _, node := range accountMap {
		if node.Parent == nil {
			accounts = append(accounts, node.Account)

			for _, child := range node.Children {
				accounts = append(accounts, child.Account)
			}
			break
		}
	}

	return accounts
}
