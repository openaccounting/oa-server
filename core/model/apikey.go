package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
)

type ApiKeyInterface interface {
	CreateApiKey(*types.ApiKey) error
	UpdateApiKey(*types.ApiKey) error
	DeleteApiKey(string, string) error
	GetApiKeys(string) ([]*types.ApiKey, error)
}

func (model *Model) CreateApiKey(key *types.ApiKey) error {
	if key.Id == "" {
		return errors.New("id required")
	}

	return model.db.InsertApiKey(key)
}

func (model *Model) UpdateApiKey(key *types.ApiKey) error {
	if key.Id == "" {
		return errors.New("id required")
	}

	return model.db.UpdateApiKey(key)
}

func (model *Model) DeleteApiKey(id string, userId string) error {
	return model.db.DeleteApiKey(id, userId)
}

func (model *Model) GetApiKeys(userId string) ([]*types.ApiKey, error) {
	return model.db.GetApiKeys(userId)
}
