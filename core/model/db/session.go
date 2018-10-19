package db

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"time"
)

type SessionInterface interface {
	InsertSession(*types.Session) error
	DeleteSession(string, string) error
	UpdateSessionActivity(string) error
}

func (db *DB) InsertSession(session *types.Session) error {
	session.Inserted = time.Now()
	session.Updated = session.Inserted

	query := "INSERT INTO session(id,inserted,updated,userId) VALUES(UNHEX(?),?,?,UNHEX(?))"
	res, err := db.Exec(
		query,
		session.Id,
		util.TimeToMs(session.Inserted),
		util.TimeToMs(session.Updated),
		session.UserId,
	)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowCnt < 1 {
		return errors.New("Unable to insert session into db")
	}

	return nil
}

func (db *DB) DeleteSession(id string, userId string) error {
	query := "UPDATE session SET `terminated` = ? WHERE id = UNHEX(?) AND userId = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(time.Now()),
		id,
		userId,
	)

	return err
}

func (db *DB) UpdateSessionActivity(id string) error {
	query := "UPDATE session SET updated = ? WHERE id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(time.Now()),
		id,
	)

	return err
}
