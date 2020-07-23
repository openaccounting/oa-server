package db

import (
	"database/sql"
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"strings"
	"time"
)

const userFields = "LOWER(HEX(u.id)),u.inserted,u.updated,u.firstName,u.lastName,u.email,u.passwordHash,u.agreeToTerms,u.passwordReset,u.emailVerified,u.emailVerifyCode,u.signupSource"

type UserInterface interface {
	InsertUser(*types.User) error
	VerifyUser(string) error
	UpdateUser(*types.User) error
	UpdateUserResetPassword(*types.User) error
	GetVerifiedUserByEmail(string) (*types.User, error)
	GetUserByActiveSession(string) (*types.User, error)
	GetUserByApiKey(string) (*types.User, error)
	GetUserByResetCode(string) (*types.User, error)
	GetUserByEmailVerifyCode(string) (*types.User, error)
	GetOrgAdmins(string) ([]*types.User, error)
}

func (db *DB) InsertUser(user *types.User) error {
	user.Inserted = time.Now()
	user.Updated = user.Inserted
	user.PasswordReset = ""

	query := "INSERT INTO user(id,inserted,updated,firstName,lastName,email,passwordHash,agreeToTerms,passwordReset,emailVerified,emailVerifyCode,signupSource) VALUES(UNHEX(?),?,?,?,?,?,?,?,?,?,?,?)"
	res, err := db.Exec(
		query,
		user.Id,
		util.TimeToMs(user.Inserted),
		util.TimeToMs(user.Updated),
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.AgreeToTerms,
		user.PasswordReset,
		user.EmailVerified,
		user.EmailVerifyCode,
		user.SignupSource,
	)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowCnt < 1 {
		return errors.New("Unable to insert user into db")
	}

	return nil
}

func (db *DB) VerifyUser(code string) error {
	query := "UPDATE user SET updated = ?, emailVerified = 1 WHERE emailVerifyCode = ?"
	res, err := db.Exec(
		query,
		util.TimeToMs(time.Now()),
		code,
	)

	count, err := res.RowsAffected()

	if err != nil {
		return nil
	}

	if count == 0 {
		return errors.New("Invalid code")
	}

	return nil
}

func (db *DB) UpdateUser(user *types.User) error {
	user.Updated = time.Now()

	query := "UPDATE user SET updated = ?, passwordHash = ?, passwordReset = ? WHERE id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(user.Updated),
		user.PasswordHash,
		"",
		user.Id,
	)

	return err
}

func (db *DB) UpdateUserResetPassword(user *types.User) error {
	user.Updated = time.Now()

	query := "UPDATE user SET updated = ?, passwordReset = ? WHERE id = UNHEX(?)"
	_, err := db.Exec(
		query,
		util.TimeToMs(user.Updated),
		user.PasswordReset,
		user.Id,
	)

	return err
}

func (db *DB) GetVerifiedUserByEmail(email string) (*types.User, error) {
	query := "SELECT  " + userFields + " FROM user u WHERE email = TRIM(LOWER(?)) AND emailVerified = 1"

	row := db.QueryRow(query, strings.TrimSpace(strings.ToLower(email)))
	u, err := db.unmarshalUser(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DB) GetUserByActiveSession(sessionId string) (*types.User, error) {
	qSelect := "SELECT " + userFields
	qFrom := " FROM user u"
	qJoin := " JOIN session s ON s.userId = u.id"
	qWhere := " WHERE s.terminated IS NULL AND s.id = UNHEX(?)"

	query := qSelect + qFrom + qJoin + qWhere

	row := db.QueryRow(query, sessionId)
	u, err := db.unmarshalUser(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DB) GetUserByApiKey(keyId string) (*types.User, error) {
	qSelect := "SELECT " + userFields
	qFrom := " FROM user u"
	qJoin := " JOIN apikey a ON a.userId = u.id"
	qWhere := " WHERE a.deleted IS NULL AND a.id = UNHEX(?)"

	query := qSelect + qFrom + qJoin + qWhere

	row := db.QueryRow(query, keyId)
	u, err := db.unmarshalUser(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DB) GetUserByResetCode(code string) (*types.User, error) {
	qSelect := "SELECT " + userFields
	qFrom := " FROM user u"
	qWhere := " WHERE u.passwordReset = ?"

	query := qSelect + qFrom + qWhere

	row := db.QueryRow(query, code)
	u, err := db.unmarshalUser(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DB) GetUserByEmailVerifyCode(code string) (*types.User, error) {
	// only allow this for 3 days
	minInserted := (time.Now().UnixNano() / 1000000) - (3 * 24 * 60 * 60 * 1000)
	qSelect := "SELECT " + userFields
	qFrom := " FROM user u"
	qWhere := " WHERE u.emailVerifyCode = ? AND inserted > ?"

	query := qSelect + qFrom + qWhere

	row := db.QueryRow(query, code, minInserted)
	u, err := db.unmarshalUser(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DB) GetOrgAdmins(orgId string) ([]*types.User, error) {
	qSelect := "SELECT " + userFields
	qFrom := " FROM user u"
	qJoin := " JOIN userorg uo ON uo.userId = u.id"
	qWhere := " WHERE uo.admin = true AND uo.orgId = UNHEX(?)"

	query := qSelect + qFrom + qJoin + qWhere

	rows, err := db.Query(query, orgId)

	if err != nil {
		return nil, err
	}

	return db.unmarshalUsers(rows)
}

func (db *DB) unmarshalUser(row *sql.Row) (*types.User, error) {
	u := new(types.User)

	var inserted int64
	var updated int64

	err := row.Scan(
		&u.Id,
		&inserted,
		&updated,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.PasswordHash,
		&u.AgreeToTerms,
		&u.PasswordReset,
		&u.EmailVerified,
		&u.EmailVerifyCode,
		&u.SignupSource,
	)

	if err != nil {
		return nil, err
	}

	u.Inserted = util.MsToTime(inserted)
	u.Updated = util.MsToTime(updated)

	return u, nil
}

func (db *DB) unmarshalUsers(rows *sql.Rows) ([]*types.User, error) {
	defer rows.Close()

	users := make([]*types.User, 0)

	for rows.Next() {
		u := new(types.User)
		var inserted int64
		var updated int64

		err := rows.Scan(
			&u.Id,
			&inserted,
			&updated,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.PasswordHash,
			&u.AgreeToTerms,
			&u.PasswordReset,
			&u.EmailVerified,
			&u.EmailVerifyCode,
			&u.SignupSource,
		)

		if err != nil {
			return nil, err
		}

		u.Inserted = util.MsToTime(inserted)
		u.Updated = util.MsToTime(updated)

		users = append(users, u)
	}

	err := rows.Err()

	return users, err
}
