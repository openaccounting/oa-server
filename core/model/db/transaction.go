package db

import (
	"database/sql"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"strconv"
	"strings"
	"time"
)

const txFields = "LOWER(HEX(id)),LOWER(HEX(orgId)),LOWER(HEX(userId)),date,inserted,updated,description,data,deleted"
const splitFields = "id,LOWER(HEX(transactionId)),LOWER(HEX(accountId)),date,inserted,updated,amount,nativeAmount,deleted"

type TransactionInterface interface {
	InsertTransaction(*types.Transaction) error
	GetTransactionById(string) (*types.Transaction, error)
	GetTransactionsByAccount(string, *types.QueryOptions) ([]*types.Transaction, error)
	GetTransactionsByOrg(string, *types.QueryOptions, []string) ([]*types.Transaction, error)
	DeleteTransaction(string) error
	DeleteAndInsertTransaction(string, *types.Transaction) error
}

func (db *DB) InsertTransaction(transaction *types.Transaction) (err error) {
	// Save to db
	dbTx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			dbTx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			dbTx.Rollback()
		} else {
			err = dbTx.Commit()
		}
	}()

	// save tx
	query1 := "INSERT INTO transaction(id,orgId,userId,date,inserted,updated,description,data) VALUES(UNHEX(?),UNHEX(?),UNHEX(?),?,?,?,?,?)"

	_, err = dbTx.Exec(
		query1,
		transaction.Id,
		transaction.OrgId,
		transaction.UserId,
		util.TimeToMs(transaction.Date),
		util.TimeToMs(transaction.Inserted),
		util.TimeToMs(transaction.Updated),
		transaction.Description,
		transaction.Data,
	)

	if err != nil {
		return
	}

	// save splits
	for _, split := range transaction.Splits {
		query := "INSERT INTO split(transactionId,accountId,date,inserted,updated,amount,nativeAmount) VALUES (UNHEX(?),UNHEX(?),?,?,?,?,?)"

		_, err = dbTx.Exec(
			query,
			transaction.Id,
			split.AccountId,
			util.TimeToMs(transaction.Date),
			util.TimeToMs(transaction.Inserted),
			util.TimeToMs(transaction.Updated),
			split.Amount,
			split.NativeAmount)

		if err != nil {
			return
		}
	}

	return
}

func (db *DB) GetTransactionById(id string) (*types.Transaction, error) {
	row := db.QueryRow("SELECT "+txFields+" FROM transaction WHERE id = UNHEX(?)", id)

	t, err := db.unmarshalTransaction(row)

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT "+splitFields+" FROM split WHERE transactionId = UNHEX(?) ORDER BY id", t.Id)

	if err != nil {
		return nil, err
	}

	t.Splits, err = db.unmarshalSplits(rows)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *DB) GetTransactionsByAccount(accountId string, options *types.QueryOptions) ([]*types.Transaction, error) {
	query := "SELECT LOWER(HEX(s.transactionId)) FROM split s"

	if options.DescriptionStartsWith != "" {
		query = query + " JOIN transaction t ON t.id = s.transactionId"
	}

	query = query + " WHERE s.accountId = UNHEX(?)"

	query = db.addOptionsToQuery(query, options)

	rows, err := db.Query(query, accountId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []string

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, "UNHEX(\""+id+"\")")
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return make([]*types.Transaction, 0), nil
	}

	query = "SELECT " + txFields + " FROM transaction WHERE id IN (" + strings.Join(ids, ",") + ")"

	query = db.addSortToQuery(query, options)

	rows, err = db.Query(query)

	if err != nil {
		return nil, err
	}

	transactions, err := db.unmarshalTransactions(rows)

	if err != nil {
		return nil, err
	}

	transactionMap := make(map[string]*types.Transaction)

	for _, t := range transactions {
		transactionMap[t.Id] = t
	}

	rows, err = db.Query("SELECT " + splitFields + " FROM split WHERE transactionId IN (" + strings.Join(ids, ",") + ") ORDER BY id")

	if err != nil {
		return nil, err
	}

	splits, err := db.unmarshalSplits(rows)

	if err != nil {
		return nil, err
	}

	for _, s := range splits {
		transaction := transactionMap[s.TransactionId]
		transaction.Splits = append(transaction.Splits, s)
	}

	return transactions, nil
}

func (db *DB) GetTransactionsByOrg(orgId string, options *types.QueryOptions, accountIds []string) ([]*types.Transaction, error) {
	if len(accountIds) == 0 {
		return make([]*types.Transaction, 0), nil
	}

	for i, accountId := range accountIds {
		accountIds[i] = "UNHEX(\"" + accountId + "\")"
	}

	query := "SELECT DISTINCT LOWER(HEX(s.transactionId)),s.date,s.inserted,s.updated FROM split s"

	if options.DescriptionStartsWith != "" {
		query = query + " JOIN transaction t ON t.id = s.transactionId"
	}

	query = query + " WHERE s.accountId IN (" + strings.Join(accountIds, ",") + ")"

	query = db.addOptionsToQuery(query, options)

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ids := []string{}

	for rows.Next() {
		var id string
		var date int64
		var inserted int64
		var updated int64
		err = rows.Scan(&id, &date, &inserted, &updated)

		if err != nil {
			return nil, err
		}

		ids = append(ids, "UNHEX(\""+id+"\")")
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return make([]*types.Transaction, 0), nil
	}

	query = "SELECT " + txFields + " FROM transaction WHERE id IN (" + strings.Join(ids, ",") + ")"

	query = db.addSortToQuery(query, options)

	rows, err = db.Query(query)

	if err != nil {
		return nil, err
	}

	transactions, err := db.unmarshalTransactions(rows)

	if err != nil {
		return nil, err
	}

	transactionMap := make(map[string]*types.Transaction)

	for _, t := range transactions {
		transactionMap[t.Id] = t
	}

	rows, err = db.Query("SELECT " + splitFields + " FROM split WHERE transactionId IN (" + strings.Join(ids, ",") + ") ORDER BY id")

	if err != nil {
		return nil, err
	}

	splits, err := db.unmarshalSplits(rows)

	if err != nil {
		return nil, err
	}

	for _, s := range splits {
		transaction := transactionMap[s.TransactionId]
		transaction.Splits = append(transaction.Splits, s)
	}

	return transactions, nil
}

func (db *DB) DeleteTransaction(id string) (err error) {
	dbTx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			dbTx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			dbTx.Rollback()
		} else {
			err = dbTx.Commit()
		}
	}()

	updatedTime := util.TimeToMs(time.Now())

	// mark splits as deleted

	query1 := "UPDATE split SET updated = ?, deleted = true WHERE transactionId = UNHEX(?)"

	_, err = dbTx.Exec(
		query1,
		updatedTime,
		id,
	)

	if err != nil {
		return
	}

	// mark transaction as deleted

	query2 := "UPDATE transaction SET updated = ?, deleted = true WHERE id = UNHEX(?)"

	_, err = dbTx.Exec(
		query2,
		updatedTime,
		id,
	)

	if err != nil {
		return
	}

	return
}

func (db *DB) DeleteAndInsertTransaction(oldId string, transaction *types.Transaction) (err error) {
	// Save to db
	dbTx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			dbTx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			dbTx.Rollback()
		} else {
			err = dbTx.Commit()
		}
	}()

	updatedTime := util.TimeToMs(transaction.Updated)

	// mark splits as deleted

	query1 := "UPDATE split SET updated = ?, deleted = true WHERE transactionId = UNHEX(?)"

	_, err = dbTx.Exec(
		query1,
		updatedTime,
		oldId,
	)

	if err != nil {
		return
	}

	// mark transaction as deleted

	query2 := "UPDATE transaction SET updated = ?, deleted = true WHERE id = UNHEX(?)"

	_, err = dbTx.Exec(
		query2,
		updatedTime,
		oldId,
	)

	if err != nil {
		return
	}

	// save new tx
	query3 := "INSERT INTO transaction(id,orgId,userId,date,inserted,updated,description,data) VALUES(UNHEX(?),UNHEX(?),UNHEX(?),?,?,?,?,?)"

	_, err = dbTx.Exec(
		query3,
		transaction.Id,
		transaction.OrgId,
		transaction.UserId,
		util.TimeToMs(transaction.Date),
		util.TimeToMs(transaction.Inserted),
		updatedTime,
		transaction.Description,
		transaction.Data,
	)

	if err != nil {
		return
	}

	// save splits
	for _, split := range transaction.Splits {
		query := "INSERT INTO split(transactionId,accountId,date,inserted,updated,amount,nativeAmount) VALUES (UNHEX(?),UNHEX(?),?,?,?,?,?)"

		_, err = dbTx.Exec(
			query,
			transaction.Id,
			split.AccountId,
			util.TimeToMs(transaction.Date),
			util.TimeToMs(transaction.Inserted),
			updatedTime,
			split.Amount,
			split.NativeAmount)

		if err != nil {
			return
		}
	}

	return
}

func (db *DB) unmarshalTransaction(row *sql.Row) (*types.Transaction, error) {
	t := new(types.Transaction)

	var date int64
	var inserted int64
	var updated int64

	err := row.Scan(&t.Id, &t.OrgId, &t.UserId, &date, &inserted, &updated, &t.Description, &t.Data, &t.Deleted)

	if err != nil {
		return nil, err
	}

	t.Date = util.MsToTime(date)
	t.Inserted = util.MsToTime(inserted)
	t.Updated = util.MsToTime(updated)

	return t, nil
}

func (db *DB) unmarshalTransactions(rows *sql.Rows) ([]*types.Transaction, error) {
	defer rows.Close()

	transactions := make([]*types.Transaction, 0)

	for rows.Next() {
		t := new(types.Transaction)
		var date int64
		var inserted int64
		var updated int64
		err := rows.Scan(&t.Id, &t.OrgId, &t.UserId, &date, &inserted, &updated, &t.Description, &t.Data, &t.Deleted)
		if err != nil {
			return nil, err
		}

		t.Date = util.MsToTime(date)
		t.Inserted = util.MsToTime(inserted)
		t.Updated = util.MsToTime(updated)
		transactions = append(transactions, t)
	}

	err := rows.Err()

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (db *DB) unmarshalSplits(rows *sql.Rows) ([]*types.Split, error) {
	defer rows.Close()

	splits := make([]*types.Split, 0)

	for rows.Next() {
		s := new(types.Split)
		var id int64
		var date int64
		var inserted int64
		var updated int64
		var deleted bool
		err := rows.Scan(&id, &s.TransactionId, &s.AccountId, &date, &inserted, &updated, &s.Amount, &s.NativeAmount, &deleted)
		if err != nil {
			return nil, err
		}

		splits = append(splits, s)
	}

	err := rows.Err()

	if err != nil {
		return nil, err
	}

	return splits, nil
}

func (db *DB) addOptionsToQuery(query string, options *types.QueryOptions) string {
	if options.IncludeDeleted != true {
		query += " AND s.deleted = false"
	}

	if options.SinceInserted != 0 {
		query += " AND s.inserted > " + strconv.Itoa(options.SinceInserted)
	}

	if options.SinceUpdated != 0 {
		query += " AND s.updated > " + strconv.Itoa(options.SinceUpdated)
	}

	if options.BeforeInserted != 0 {
		query += " AND s.inserted < " + strconv.Itoa(options.BeforeInserted)
	}

	if options.BeforeUpdated != 0 {
		query += " AND s.updated < " + strconv.Itoa(options.BeforeUpdated)
	}

	if options.StartDate != 0 {
		query += " AND s.date >= " + strconv.Itoa(options.StartDate)
	}

	if options.EndDate != 0 {
		query += " AND s.date < " + strconv.Itoa(options.EndDate)
	}

	if options.DescriptionStartsWith != "" {
		query += " AND t.description LIKE '" + db.Escape(options.DescriptionStartsWith) + "%'"
	}

	if options.Sort == "updated-asc" {
		query += " ORDER BY s.updated ASC"
	} else {
		query += " ORDER BY s.date DESC, s.inserted DESC"
	}

	if options.Limit != 0 && options.Skip != 0 {
		query += " LIMIT " + strconv.Itoa(options.Skip) + ", " + strconv.Itoa(options.Limit)
	} else if options.Limit != 0 {
		query += " LIMIT " + strconv.Itoa(options.Limit)
	}

	return query
}

func (db *DB) addSortToQuery(query string, options *types.QueryOptions) string {
	if options.Sort == "updated-asc" {
		query += " ORDER BY updated ASC"
	} else {
		query += " ORDER BY date DESC, inserted DESC"
	}

	return query
}
