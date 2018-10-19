package db

import (
	"database/sql"
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"time"
)

type OrgInterface interface {
	CreateOrg(*types.Org, string, []*types.Account) error
	UpdateOrg(*types.Org) error
	GetOrg(string, string) (*types.Org, error)
	GetOrgs(string) ([]*types.Org, error)
	GetOrgUserIds(string) ([]string, error)
	InsertInvite(*types.Invite) error
	AcceptInvite(*types.Invite, string) error
	GetInvites(string) ([]*types.Invite, error)
	GetInvite(string) (*types.Invite, error)
	DeleteInvite(string) error
}

const orgFields = "LOWER(HEX(o.id)),o.inserted,o.updated,o.name,o.currency,o.`precision`"
const inviteFields = "i.id,LOWER(HEX(i.orgId)),i.inserted,i.updated,i.email,i.accepted"

func (db *DB) CreateOrg(org *types.Org, userId string, accounts []*types.Account) (err error) {
	tx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	org.Inserted = time.Now()
	org.Updated = org.Inserted

	// create org
	query1 := "INSERT INTO org(id,inserted,updated,name,currency,`precision`) VALUES(UNHEX(?),?,?,?,?,?)"

	res, err := tx.Exec(
		query1,
		org.Id,
		util.TimeToMs(org.Inserted),
		util.TimeToMs(org.Updated),
		org.Name,
		org.Currency,
		org.Precision,
	)

	if err != nil {
		return
	}

	// associate user with org
	query2 := "INSERT INTO userorg(userId,orgId,admin) VALUES(UNHEX(?),UNHEX(?), 1)"

	res, err = tx.Exec(query2, userId, org.Id)

	if err != nil {
		return
	}

	_, err = res.LastInsertId()

	if err != nil {
		return
	}

	// create Accounts: Root, Assets, Liabilities, Equity, Income, Expenses

	for _, account := range accounts {

		query := "INSERT INTO account(id,orgId,inserted,updated,name,parent,currency,`precision`,debitBalance) VALUES (UNHEX(?),UNHEX(?),?,?,?,UNHEX(?),?,?,?)"

		if _, err = tx.Exec(
			query,
			account.Id,
			org.Id,
			util.TimeToMs(org.Inserted),
			util.TimeToMs(org.Updated),
			account.Name,
			account.Parent,
			account.Currency,
			account.Precision,
			account.DebitBalance,
		); err != nil {
			return
		}
	}

	permissionId, err := util.NewGuid()

	if err != nil {
		return
	}

	// Grant root permission to user

	query3 := "INSERT INTO permission (id,userId,orgId,accountId,type,inserted,updated) VALUES(UNHEX(?),UNHEX(?),UNHEX(?),UNHEX(?),?,?,?)"

	_, err = tx.Exec(
		query3,
		permissionId,
		userId,
		org.Id,
		accounts[0].Id,
		0,
		util.TimeToMs(org.Inserted),
		util.TimeToMs(org.Updated),
	)

	return
}

func (db *DB) UpdateOrg(org *types.Org) error {
	org.Updated = time.Now()

	query := "UPDATE org SET updated = ?, name = ? WHERE id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(org.Updated),
		org.Name,
		org.Id,
	)

	return err
}

func (db *DB) GetOrg(orgId string, userId string) (*types.Org, error) {
	var o types.Org
	var inserted int64
	var updated int64

	err := db.QueryRow("SELECT "+orgFields+" FROM org o JOIN userorg ON userorg.orgId = o.id WHERE o.id = UNHEX(?) AND userorg.userId = UNHEX(?)", orgId, userId).
		Scan(&o.Id, &inserted, &updated, &o.Name, &o.Currency, &o.Precision)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Org not found")
	case err != nil:
		return nil, err
	default:
		o.Inserted = util.MsToTime(inserted)
		o.Updated = util.MsToTime(updated)
		return &o, nil
	}
}

func (db *DB) GetOrgs(userId string) ([]*types.Org, error) {
	rows, err := db.Query("SELECT "+orgFields+" from org o JOIN userorg ON userorg.orgId = o.id WHERE userorg.userId = UNHEX(?)", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orgs := make([]*types.Org, 0)

	for rows.Next() {
		o := new(types.Org)
		var inserted int64
		var updated int64

		err = rows.Scan(&o.Id, &inserted, &updated, &o.Name, &o.Currency, &o.Precision)
		if err != nil {
			return nil, err
		}

		o.Inserted = util.MsToTime(inserted)
		o.Updated = util.MsToTime(updated)

		orgs = append(orgs, o)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (db *DB) GetOrgUserIds(orgId string) ([]string, error) {
	rows, err := db.Query("SELECT LOWER(HEX(userId)) FROM userorg WHERE orgId = UNHEX(?)", orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userIds := make([]string, 0)

	for rows.Next() {
		var userId string
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}

		userIds = append(userIds, userId)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userIds, nil
}

func (db *DB) InsertInvite(invite *types.Invite) error {
	invite.Inserted = time.Now()
	invite.Updated = invite.Inserted

	query := "INSERT INTO invite(id,orgId,inserted,updated,email,accepted) VALUES(?,UNHEX(?),?,?,?,?)"
	_, err := db.Exec(
		query,
		invite.Id,
		invite.OrgId,
		util.TimeToMs(invite.Inserted),
		util.TimeToMs(invite.Updated),
		invite.Email,
		false,
	)

	return err
}

func (db *DB) AcceptInvite(invite *types.Invite, userId string) error {
	invite.Updated = time.Now()

	// Get root account for permission
	rootAccount, err := db.GetRootAccount(invite.OrgId)

	if err != nil {
		return err
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// associate user with org
	query1 := "INSERT INTO userorg(userId,orgId,admin) VALUES(UNHEX(?),UNHEX(?), 0)"

	_, err = tx.Exec(query1, userId, invite.OrgId)

	if err != nil {
		return err
	}

	query2 := "UPDATE invite SET accepted = 1, updated = ? WHERE id = ?"

	_, err = tx.Exec(query2, util.TimeToMs(invite.Updated), invite.Id)

	// Grant root permission to user

	permissionId, err := util.NewGuid()

	if err != nil {
		return err
	}

	query3 := "INSERT INTO permission (id,userId,orgId,accountId,type,inserted,updated) VALUES(UNHEX(?),UNHEX(?),UNHEX(?),UNHEX(?),?,?,?)"

	_, err = tx.Exec(
		query3,
		permissionId,
		userId,
		invite.OrgId,
		rootAccount.Id,
		0,
		util.TimeToMs(invite.Updated),
		util.TimeToMs(invite.Updated),
	)

	return err
}

func (db *DB) GetInvites(orgId string) ([]*types.Invite, error) {
	// don't include expired invoices
	cutoff := util.TimeToMs(time.Now()) - 7*24*60*60*1000

	rows, err := db.Query("SELECT "+inviteFields+" FROM invite i WHERE orgId = UNHEX(?) AND inserted > ?", orgId, cutoff)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	invites := make([]*types.Invite, 0)

	for rows.Next() {
		i := new(types.Invite)
		var inserted int64
		var updated int64

		err = rows.Scan(&i.Id, &i.OrgId, &inserted, &updated, &i.Email, &i.Accepted)
		if err != nil {
			return nil, err
		}

		i.Inserted = util.MsToTime(inserted)
		i.Updated = util.MsToTime(updated)

		invites = append(invites, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (db *DB) GetInvite(id string) (*types.Invite, error) {
	var i types.Invite
	var inserted int64
	var updated int64

	err := db.QueryRow("SELECT "+inviteFields+" FROM invite i WHERE i.id = ?", id).
		Scan(&i.Id, &i.OrgId, &inserted, &updated, &i.Email, &i.Accepted)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Invite not found")
	case err != nil:
		return nil, err
	default:
		i.Inserted = util.MsToTime(inserted)
		i.Updated = util.MsToTime(updated)
		return &i, nil
	}
}

func (db *DB) DeleteInvite(id string) error {
	query := "DELETE FROM invite WHERE id = ?"
	_, err := db.Exec(
		query,
		id,
	)

	return err
}
