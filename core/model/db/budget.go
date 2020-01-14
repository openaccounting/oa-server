package db

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"time"
)

type BudgetInterface interface {
	GetBudget(string) (*types.Budget, error)
	InsertAndReplaceBudget(*types.Budget) error
	DeleteBudget(string) error
}

const budgetFields = "LOWER(HEX(accountId)),inserted,amount"

func (db *DB) GetBudget(orgId string) (*types.Budget, error) {
	var budget types.Budget
	var inserted int64

	rows, err := db.Query("SELECT "+budgetFields+" FROM budgetitem WHERE orgId = UNHEX(?) ORDER BY HEX(accountId)", orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]*types.BudgetItem, 0)

	for rows.Next() {
		i := new(types.BudgetItem)
		err := rows.Scan(&i.AccountId, &inserted, &i.Amount)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("Budget not found")
	}

	budget.OrgId = orgId
	budget.Inserted = util.MsToTime(inserted)
	budget.Items = items

	return &budget, nil
}

func (db *DB) InsertAndReplaceBudget(budget *types.Budget) (err error) {
	budget.Inserted = time.Now()

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

	// delete previous budget
	query1 := "DELETE FROM budgetitem WHERE orgId = UNHEX(?)"

	_, err = dbTx.Exec(
		query1,
		budget.OrgId,
	)

	if err != nil {
		return
	}

	// save items
	for _, item := range budget.Items {
		query := "INSERT INTO budgetitem(orgId,accountId,inserted,amount) VALUES (UNHEX(?),UNHEX(?),?,?)"

		_, err = dbTx.Exec(
			query,
			budget.OrgId,
			item.AccountId,
			util.TimeToMs(budget.Inserted),
			item.Amount)

		if err != nil {
			return
		}
	}

	return
}

func (db *DB) DeleteBudget(orgId string) error {
	query := "DELETE FROM budgetitem WHERE orgId = UNHEX(?)"

	_, err := db.Exec(query, orgId)

	return err
}
