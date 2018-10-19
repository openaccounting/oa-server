package db

import (
	"database/sql"
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"time"
)

type PriceInterface interface {
	InsertPrice(*types.Price) error
	GetPriceById(string) (*types.Price, error)
	DeletePrice(string) error
	GetPricesNearestInTime(string, time.Time) ([]*types.Price, error)
	GetPricesByCurrency(string, string) ([]*types.Price, error)
}

const priceFields = "LOWER(HEX(p.id)),LOWER(HEX(p.orgId)),p.currency,p.date,p.inserted,p.updated,p.price"

func (db *DB) InsertPrice(price *types.Price) error {
	price.Inserted = time.Now()
	price.Updated = price.Inserted

	if price.Date.IsZero() {
		price.Date = price.Inserted
	}

	query := "INSERT INTO price(id,orgId,currency,date,inserted,updated,price) VALUES(UNHEX(?),UNHEX(?),?,?,?,?,?)"
	_, err := db.Exec(
		query,
		price.Id,
		price.OrgId,
		price.Currency,
		util.TimeToMs(price.Date),
		util.TimeToMs(price.Inserted),
		util.TimeToMs(price.Updated),
		price.Price,
	)

	return err
}

func (db *DB) GetPriceById(id string) (*types.Price, error) {
	var p types.Price
	var date int64
	var inserted int64
	var updated int64

	err := db.QueryRow("SELECT "+priceFields+" FROM price p WHERE id = UNHEX(?)", id).
		Scan(&p.Id, &p.OrgId, &p.Currency, &date, &inserted, &updated, &p.Price)

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("Price not found")
	case err != nil:
		return nil, err
	default:
		p.Date = util.MsToTime(date)
		p.Inserted = util.MsToTime(inserted)
		p.Updated = util.MsToTime(updated)
		return &p, nil
	}
}

func (db *DB) DeletePrice(id string) error {
	query := "DELETE FROM price WHERE id = UNHEX(?)"

	_, err := db.Exec(query, id)

	return err
}

func (db *DB) GetPricesNearestInTime(orgId string, date time.Time) ([]*types.Price, error) {
	qSelect := "SELECT " + priceFields
	qFrom := " FROM price p"
	qJoin := " LEFT OUTER JOIN price p2 ON p.currency = p2.currency AND p.orgId = p2.orgId AND ABS(CAST(p.date AS SIGNED) - ?) > ABS(CAST(p2.date AS SIGNED) - ?)"
	qWhere := " WHERE p2.id IS NULL AND p.orgId = UNHEX(?)"

	query := qSelect + qFrom + qJoin + qWhere

	rows, err := db.Query(query, util.TimeToMs(date), util.TimeToMs(date), orgId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	prices := make([]*types.Price, 0)

	for rows.Next() {
		var date int64
		var inserted int64
		var updated int64
		p := new(types.Price)
		err = rows.Scan(&p.Id, &p.OrgId, &p.Currency, &date, &inserted, &updated, &p.Price)
		if err != nil {
			return nil, err
		}

		p.Date = util.MsToTime(date)
		p.Inserted = util.MsToTime(inserted)
		p.Updated = util.MsToTime(updated)

		prices = append(prices, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (db *DB) GetPricesByCurrency(orgId string, currency string) ([]*types.Price, error) {
	qSelect := "SELECT " + priceFields
	qFrom := " FROM price p"
	qWhere := " WHERE p.orgId = UNHEX(?) AND p.currency = ?"
	pOrder := " ORDER BY date ASC"

	query := qSelect + qFrom + qWhere + pOrder

	rows, err := db.Query(query, orgId, currency)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	prices := make([]*types.Price, 0)

	for rows.Next() {
		var date int64
		var inserted int64
		var updated int64
		p := new(types.Price)
		err = rows.Scan(&p.Id, &p.OrgId, &p.Currency, &date, &inserted, &updated, &p.Price)
		if err != nil {
			return nil, err
		}

		p.Date = util.MsToTime(date)
		p.Inserted = util.MsToTime(inserted)
		p.Updated = util.MsToTime(updated)

		prices = append(prices, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return prices, nil
}
