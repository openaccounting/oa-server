package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
)

type SessionInterface interface {
	CreateSession(*types.Session) error
	DeleteSession(string, string) error
}

func (model *Model) CreateSession(session *types.Session) error {
	if session.Id == "" {
		return errors.New("id required")
	}

	return model.db.InsertSession(session)
}

func (model *Model) DeleteSession(id string, userId string) error {
	return model.db.DeleteSession(id, userId)
}
