package model

import (
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
)

var Instance Interface

type Model struct {
	db     db.Datastore
	bcrypt util.Bcrypt
	config types.Config
}

type Interface interface {
	UserInterface
	OrgInterface
	AccountInterface
	TransactionInterface
	PriceInterface
	SessionInterface
	ApiKeyInterface
}

func NewModel(db db.Datastore, bcrypt util.Bcrypt, config types.Config) *Model {
	model := &Model{db: db, bcrypt: bcrypt, config: config}
	Instance = model
	return model
}
