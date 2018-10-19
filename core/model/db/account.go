package db

import (
	"database/sql"
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"math"
	"strings"
	"time"
)

const emptyAccountId = "00000000000000000000000000000000"

type AccountInterface interface {
	InsertAccount(account *types.Account) error
	UpdateAccount(account *types.Account) error
	GetAccount(string) (*types.Account, error)
	GetAccountsByOrgId(orgId string) ([]*types.Account, error)
	GetPermissionedAccountIds(string, string, string) ([]string, error)
	GetSplitCountByAccountId(id string) (int64, error)
	GetChildCountByAccountId(id string) (int64, error)
	DeleteAccount(id string) error
	AddBalances([]*types.Account, time.Time) error
	AddNativeBalancesCost([]*types.Account, time.Time) error
	AddNativeBalancesNearestInTime([]*types.Account, time.Time) error
	AddBalance(*types.Account, time.Time) error
	AddNativeBalanceCost(*types.Account, time.Time) error
	AddNativeBalanceNearestInTime(*types.Account, time.Time) error
	GetRootAccount(string) (*types.Account, error)
}

func (db *DB) InsertAccount(account *types.Account) error {
	account.Inserted = time.Now()
	account.Updated = account.Inserted

	query := "INSERT INTO account(id,orgId,inserted,updated,name,parent,currency,`precision`,debitBalance) VALUES(UNHEX(?),UNHEX(?),?,?,?,UNHEX(?),?,?,?)"
	_, err := db.Exec(
		query,
		account.Id,
		account.OrgId,
		util.TimeToMs(account.Inserted),
		util.TimeToMs(account.Updated),
		account.Name,
		account.Parent,
		account.Currency,
		account.Precision,
		account.DebitBalance)

	return err
}

func (db *DB) UpdateAccount(account *types.Account) error {
	account.Updated = time.Now()

	query := "UPDATE account SET updated = ?, name = ?, parent = UNHEX(?), currency = ?, `precision` = ?, debitBalance = ? WHERE id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(account.Updated),
		account.Name,
		account.Parent,
		account.Currency,
		account.Precision,
		account.DebitBalance,
		account.Id)

	return err
}

func (db *DB) GetAccount(id string) (*types.Account, error) {
	a := types.Account{}
	var inserted int64
	var updated int64

	err := db.QueryRow("SELECT LOWER(HEX(id)),LOWER(HEX(orgId)),inserted,updated,name,LOWER(HEX(parent)),currency,`precision`,debitBalance FROM account WHERE id = UNHEX(?)", id).
		Scan(&a.Id, &a.OrgId, &inserted, &updated, &a.Name, &a.Parent, &a.Currency, &a.Precision, &a.DebitBalance)

	if a.Parent == emptyAccountId {
		a.Parent = ""
	}

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Account not found")
	case err != nil:
		return nil, err
	default:
		a.Inserted = util.MsToTime(inserted)
		a.Updated = util.MsToTime(updated)
		return &a, nil
	}
}

func (db *DB) GetAccountsByOrgId(orgId string) ([]*types.Account, error) {
	rows, err := db.Query("SELECT LOWER(HEX(id)),LOWER(HEX(orgId)),inserted,updated,name,LOWER(HEX(parent)),currency,`precision`,debitBalance FROM account WHERE orgId = UNHEX(?)", orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := make([]*types.Account, 0)

	for rows.Next() {
		a := new(types.Account)
		var inserted int64
		var updated int64

		err = rows.Scan(&a.Id, &a.OrgId, &inserted, &updated, &a.Name, &a.Parent, &a.Currency, &a.Precision, &a.DebitBalance)
		if err != nil {
			return nil, err
		}

		if a.Parent == emptyAccountId {
			a.Parent = ""
		}

		a.Inserted = util.MsToTime(inserted)
		a.Updated = util.MsToTime(updated)

		accounts = append(accounts, a)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (db *DB) GetPermissionedAccountIds(orgId string, userId string, tokenId string) ([]string, error) {
	// Get user permissions
	// TODO incorporate tokens
	rows, err := db.Query("SELECT LOWER(HEX(accountId)) FROM permission WHERE orgId = UNHEX(?) AND userId = UNHEX(?)", orgId, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissionedAccounts []string

	var id string

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		permissionedAccounts = append(permissionedAccounts, id)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return permissionedAccounts, nil
}

func (db *DB) GetSplitCountByAccountId(id string) (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM split WHERE deleted = false AND accountId = UNHEX(?)"

	err := db.QueryRow(query, id).Scan(&count)

	return count, err
}

func (db *DB) GetChildCountByAccountId(id string) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM account WHERE parent = UNHEX(?)"

	err := db.QueryRow(query, id).Scan(&count)

	return count, err
}

func (db *DB) DeleteAccount(id string) error {
	query := "DELETE FROM account WHERE id = UNHEX(?)"

	_, err := db.Exec(query, id)

	return err
}

func (db *DB) AddBalances(accounts []*types.Account, date time.Time) error {
	// TODO optimize
	ids := make([]string, len(accounts))

	for i, account := range accounts {
		ids[i] = "UNHEX(\"" + account.Id + "\")"
	}

	balanceMap := make(map[string]*int64)

	query := "SELECT LOWER(HEX(accountId)), SUM(amount) FROM split WHERE deleted = false AND accountId IN (" +
		strings.Join(ids, ",") + ")" +
		" AND date < ? GROUP BY accountId"

	rows, err := db.Query(query, util.TimeToMs(date))

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var balance int64
		err := rows.Scan(&id, &balance)
		if err != nil {
			return err
		}

		balanceMap[id] = &balance
	}

	err = rows.Err()

	if err != nil {
		return err
	}

	for _, account := range accounts {
		account.Balance = balanceMap[account.Id]
	}

	return nil
}

func (db *DB) AddNativeBalancesCost(accounts []*types.Account, date time.Time) error {
	// TODO optimize
	ids := make([]string, len(accounts))

	for i, account := range accounts {
		ids[i] = "UNHEX(\"" + account.Id + "\")"
	}

	balanceMap := make(map[string]*int64)

	query := "SELECT LOWER(HEX(accountId)), SUM(nativeAmount) FROM split WHERE deleted = false AND accountId IN (" +
		strings.Join(ids, ",") + ")" +
		" AND date < ? GROUP BY accountId"

	rows, err := db.Query(query, util.TimeToMs(date))

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var balance int64
		err := rows.Scan(&id, &balance)
		if err != nil {
			return err
		}

		balanceMap[id] = &balance
	}

	err = rows.Err()

	if err != nil {
		return err
	}

	for _, account := range accounts {
		account.NativeBalance = balanceMap[account.Id]
	}

	return nil
}

func (db *DB) AddNativeBalancesNearestInTime(accounts []*types.Account, date time.Time) error {
	// TODO Don't look up org currency every single time

	for _, account := range accounts {
		err := db.AddNativeBalanceNearestInTime(account, date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) AddBalance(account *types.Account, date time.Time) error {
	var balance sql.NullInt64

	query := "SELECT SUM(amount) FROM split WHERE deleted = false AND accountId = UNHEX(?) AND date < ?"

	err := db.QueryRow(query, account.Id, util.TimeToMs(date)).Scan(&balance)

	if err != nil {
		return err
	}

	account.Balance = &balance.Int64

	return nil
}

func (db *DB) AddNativeBalanceCost(account *types.Account, date time.Time) error {
	var nativeBalance sql.NullInt64

	query := "SELECT SUM(nativeAmount) FROM split WHERE deleted = false AND accountId = UNHEX(?) AND date < ?"

	err := db.QueryRow(query, account.Id, util.TimeToMs(date)).Scan(&nativeBalance)

	if err != nil {
		return err
	}

	account.NativeBalance = &nativeBalance.Int64

	return nil
}

func (db *DB) AddNativeBalanceNearestInTime(account *types.Account, date time.Time) error {
	var orgCurrency string
	var orgPrecision int

	query1 := "SELECT currency,`precision` FROM org WHERE id = UNHEX(?)"

	err := db.QueryRow(query1, account.OrgId).Scan(&orgCurrency, &orgPrecision)

	if err != nil {
		return err
	}

	if account.Balance == nil {
		return nil
	}

	if orgCurrency == account.Currency {
		nativeBalance := int64(*account.Balance)
		account.NativeBalance = &nativeBalance
		return nil
	}

	var tmp sql.NullInt64
	var price float64

	query2 := "SELECT ABS(CAST(date AS SIGNED) - ?) AS datediff, price FROM price WHERE currency = ? ORDER BY datediff ASC LIMIT 1"

	err = db.QueryRow(query2, util.TimeToMs(date), account.Currency).Scan(&tmp, &price)

	if err == sql.ErrNoRows {
		nativeBalance := int64(0)
		account.NativeBalance = &nativeBalance
	} else if err != nil {
		return err
	}

	precisionAdj := math.Pow(10, float64(account.Precision-orgPrecision))
	nativeBalance := int64(float64(*account.Balance) * price / precisionAdj)
	account.NativeBalance = &nativeBalance

	return nil
}

func (db *DB) GetRootAccount(orgId string) (*types.Account, error) {
	a := types.Account{}
	var inserted int64
	var updated int64

	err := db.QueryRow(
		"SELECT LOWER(HEX(id)),LOWER(HEX(orgId)),inserted,updated,name,LOWER(HEX(parent)),currency,`precision`,debitBalance FROM account WHERE orgId = UNHEX(?) AND parent = UNHEX(?)",
		orgId,
		emptyAccountId).
		Scan(&a.Id, &a.OrgId, &inserted, &updated, &a.Name, &a.Parent, &a.Currency, &a.Precision, &a.DebitBalance)

	a.Parent = ""

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Account not found")
	case err != nil:
		return nil, err
	default:
		a.Inserted = util.MsToTime(inserted)
		a.Updated = util.MsToTime(updated)
		return &a, nil
	}
}
