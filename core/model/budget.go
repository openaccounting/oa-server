package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
)

type BudgetInterface interface {
	GetBudget(string, string) (*types.Budget, error)
	CreateBudget(*types.Budget, string) error
	DeleteBudget(string, string) error
}

func (model *Model) GetBudget(orgId string, userId string) (*types.Budget, error) {
	belongs, err := model.UserBelongsToOrg(userId, orgId)

	if err != nil {
		return nil, err
	}

	if belongs == false {
		return nil, errors.New("User does not belong to org")
	}

	return model.db.GetBudget(orgId)
}

func (model *Model) CreateBudget(budget *types.Budget, userId string) error {
	belongs, err := model.UserBelongsToOrg(userId, budget.OrgId)

	if err != nil {
		return err
	}

	if belongs == false {
		return errors.New("User does not belong to org")
	}

	if budget.OrgId == "" {
		return errors.New("orgId required")
	}

	return model.db.InsertAndReplaceBudget(budget)
}

func (model *Model) DeleteBudget(orgId string, userId string) error {
	belongs, err := model.UserBelongsToOrg(userId, orgId)

	if err != nil {
		return err
	}

	if belongs == false {
		return errors.New("User does not belong to org")
	}

	return model.db.DeleteBudget(orgId)
}
