package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/ws"
	"time"
)

type PriceInterface interface {
	CreatePrice(*types.Price, string) error
	DeletePrice(string, string) error
	GetPricesNearestInTime(string, time.Time, string) ([]*types.Price, error)
	GetPricesByCurrency(string, string, string) ([]*types.Price, error)
}

func (model *Model) CreatePrice(price *types.Price, userId string) error {
	belongs, err := model.UserBelongsToOrg(userId, price.OrgId)

	if err != nil {
		return err
	}

	if belongs == false {
		return errors.New("User does not belong to org")
	}

	if price.Id == "" {
		return errors.New("id required")
	}

	if price.OrgId == "" {
		return errors.New("orgId required")
	}

	if price.Currency == "" {
		return errors.New("currency required")
	}

	err = model.db.InsertPrice(price)

	if err != nil {
		return err
	}

	// Notify web socket subscribers
	userIds, err2 := model.db.GetOrgUserIds(price.OrgId)

	if err2 == nil {
		ws.PushPrice(price, userIds, "create")
	}

	return nil
}

func (model *Model) DeletePrice(id string, userId string) error {
	// Get original price
	price, err := model.db.GetPriceById(id)

	if err != nil {
		return err
	}

	belongs, err := model.UserBelongsToOrg(userId, price.OrgId)

	if err != nil {
		return err
	}

	if belongs == false {
		return errors.New("User does not belong to org")
	}

	err = model.db.DeletePrice(id)

	if err != nil {
		return err
	}

	// Notify web socket subscribers
	// TODO only get user ids that have permission to access account
	userIds, err2 := model.db.GetOrgUserIds(price.OrgId)

	if err2 == nil {
		ws.PushPrice(price, userIds, "delete")
	}

	return nil
}

func (model *Model) GetPricesNearestInTime(orgId string, date time.Time, userId string) ([]*types.Price, error) {
	belongs, err := model.UserBelongsToOrg(userId, orgId)

	if err != nil {
		return nil, err
	}

	if belongs == false {
		return nil, errors.New("User does not belong to org")
	}

	return model.db.GetPricesNearestInTime(orgId, date)
}

func (model *Model) GetPricesByCurrency(orgId string, currency string, userId string) ([]*types.Price, error) {
	belongs, err := model.UserBelongsToOrg(userId, orgId)

	if err != nil {
		return nil, err
	}

	if belongs == false {
		return nil, errors.New("User does not belong to org")
	}

	return model.db.GetPricesByCurrency(orgId, currency)
}
