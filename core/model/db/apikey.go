package db

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"time"
)

type ApiKeyInterface interface {
	InsertApiKey(*types.ApiKey) error
	UpdateApiKey(*types.ApiKey) error
	DeleteApiKey(string, string) error
	GetApiKeys(string) ([]*types.ApiKey, error)
	UpdateApiKeyActivity(string) error
}

const apiKeyFields = "LOWER(HEX(id)),inserted,updated,LOWER(HEX(userId)),label"

func (db *DB) InsertApiKey(key *types.ApiKey) error {
	key.Inserted = time.Now()
	key.Updated = key.Inserted

	query := "INSERT INTO apikey(id,inserted,updated,userId,label) VALUES(UNHEX(?),?,?,UNHEX(?),?)"
	res, err := db.Exec(
		query,
		key.Id,
		util.TimeToMs(key.Inserted),
		util.TimeToMs(key.Updated),
		key.UserId,
		key.Label,
	)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowCnt < 1 {
		return errors.New("Unable to insert apikey into db")
	}

	return nil
}

func (db *DB) UpdateApiKey(key *types.ApiKey) error {
	key.Updated = time.Now()

	query := "UPDATE apikey SET updated = ?, label = ? WHERE deleted IS NULL AND id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(key.Updated),
		key.Label,
		key.Id,
	)

	if err != nil {
		return err
	}

	var inserted int64

	err = db.QueryRow("SELECT inserted FROM apikey WHERE id = UNHEX(?)", key.Id).Scan(&inserted)

	if err != nil {
		return err
	}

	key.Inserted = util.MsToTime(inserted)

	return nil
}

func (db *DB) DeleteApiKey(id string, userId string) error {
	query := "UPDATE apikey SET deleted = ? WHERE id = UNHEX(?) AND userId = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(time.Now()),
		id,
		userId,
	)

	return err
}

func (db *DB) GetApiKeys(userId string) ([]*types.ApiKey, error) {
	rows, err := db.Query("SELECT "+apiKeyFields+" from apikey WHERE deleted IS NULL AND userId = UNHEX(?)", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	keys := make([]*types.ApiKey, 0)

	for rows.Next() {
		k := new(types.ApiKey)
		var inserted int64
		var updated int64

		err = rows.Scan(&k.Id, &inserted, &updated, &k.UserId, &k.Label)
		if err != nil {
			return nil, err
		}

		k.Inserted = util.MsToTime(inserted)
		k.Updated = util.MsToTime(updated)

		keys = append(keys, k)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (db *DB) UpdateApiKeyActivity(id string) error {
	query := "UPDATE apikey SET updated = ? WHERE id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(time.Now()),
		id,
	)

	return err
}
