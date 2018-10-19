package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"time"
)

type OrgInterface interface {
	CreateOrg(*types.Org, string) error
	UpdateOrg(*types.Org, string) error
	GetOrg(string, string) (*types.Org, error)
	GetOrgs(string) ([]*types.Org, error)
	CreateInvite(*types.Invite, string) error
	AcceptInvite(*types.Invite, string) error
	GetInvites(string, string) ([]*types.Invite, error)
	DeleteInvite(string, string) error
}

func (model *Model) CreateOrg(org *types.Org, userId string) error {
	if org.Name == "" {
		return errors.New("name required")
	}

	if org.Currency == "" {
		return errors.New("currency required")
	}

	accounts := make([]*types.Account, 6)

	id, err := util.NewGuid()

	if err != nil {
		return err
	}

	accounts[0] = &types.Account{
		Id:           id,
		Name:         "Root",
		Parent:       "",
		Currency:     org.Currency,
		Precision:    org.Precision,
		DebitBalance: true,
	}

	id, err = util.NewGuid()

	if err != nil {
		return err
	}

	accounts[1] = &types.Account{
		Id:           id,
		Name:         "Assets",
		Parent:       accounts[0].Id,
		Currency:     org.Currency,
		Precision:    org.Precision,
		DebitBalance: true,
	}

	id, err = util.NewGuid()

	if err != nil {
		return err
	}

	accounts[2] = &types.Account{
		Id:           id,
		Name:         "Liabilities",
		Parent:       accounts[0].Id,
		Currency:     org.Currency,
		Precision:    org.Precision,
		DebitBalance: false,
	}

	id, err = util.NewGuid()

	if err != nil {
		return err
	}

	accounts[3] = &types.Account{
		Id:           id,
		Name:         "Equity",
		Parent:       accounts[0].Id,
		Currency:     org.Currency,
		Precision:    org.Precision,
		DebitBalance: false,
	}

	id, err = util.NewGuid()

	if err != nil {
		return err
	}

	accounts[4] = &types.Account{
		Id:           id,
		Name:         "Income",
		Parent:       accounts[0].Id,
		Currency:     org.Currency,
		Precision:    org.Precision,
		DebitBalance: false,
	}

	id, err = util.NewGuid()

	if err != nil {
		return err
	}

	accounts[5] = &types.Account{
		Id:           id,
		Name:         "Expenses",
		Parent:       accounts[0].Id,
		Currency:     org.Currency,
		Precision:    org.Precision,
		DebitBalance: true,
	}

	return model.db.CreateOrg(org, userId, accounts)
}

func (model *Model) UpdateOrg(org *types.Org, userId string) error {
	_, err := model.GetOrg(org.Id, userId)

	if err != nil {
		// user doesn't have access to org
		return errors.New("access denied")
	}

	if org.Name == "" {
		return errors.New("name required")
	}

	return model.db.UpdateOrg(org)
}

func (model *Model) GetOrg(orgId string, userId string) (*types.Org, error) {
	return model.db.GetOrg(orgId, userId)
}

func (model *Model) GetOrgs(userId string) ([]*types.Org, error) {
	return model.db.GetOrgs(userId)
}

func (model *Model) UserBelongsToOrg(userId string, orgId string) (bool, error) {
	orgs, err := model.GetOrgs(userId)

	if err != nil {
		return false, err
	}

	belongs := false

	for _, org := range orgs {
		if org.Id == orgId {
			belongs = true
			break
		}
	}

	return belongs, nil
}

func (model *Model) CreateInvite(invite *types.Invite, userId string) error {
	admins, err := model.db.GetOrgAdmins(invite.OrgId)

	if err != nil {
		return err
	}

	isAdmin := false

	for _, admin := range admins {
		if admin.Id == userId {
			isAdmin = true
			break
		}
	}

	if isAdmin == false {
		return errors.New("Must be org admin to invite users")
	}

	inviteId, err := util.NewInviteId()

	if err != nil {
		return err
	}

	invite.Id = inviteId

	err = model.db.InsertInvite(invite)

	if err != nil {
		return err
	}

	if invite.Email != "" {
		// TODO send email
	}

	return nil
}

func (model *Model) AcceptInvite(invite *types.Invite, userId string) error {
	if invite.Accepted != true {
		return errors.New("accepted must be true")
	}

	if invite.Id == "" {
		return errors.New("missing invite id")
	}

	// Get original invite
	original, err := model.db.GetInvite(invite.Id)

	if err != nil {
		return err
	}

	if original.Accepted == true {
		return errors.New("invite already accepted")
	}

	oneWeekAfter := original.Inserted.Add(time.Hour * 24 * 7)

	if time.Now().After(oneWeekAfter) == true {
		return errors.New("invite has expired")
	}

	invite.OrgId = original.OrgId
	invite.Email = original.Email
	invite.Inserted = original.Inserted

	return model.db.AcceptInvite(invite, userId)
}

func (model *Model) GetInvites(orgId string, userId string) ([]*types.Invite, error) {
	admins, err := model.db.GetOrgAdmins(orgId)

	if err != nil {
		return nil, err
	}

	isAdmin := false

	for _, admin := range admins {
		if admin.Id == userId {
			isAdmin = true
			break
		}
	}

	if isAdmin == false {
		return nil, errors.New("Must be org admin to invite users")
	}

	return model.db.GetInvites(orgId)
}

func (model *Model) DeleteInvite(id string, userId string) error {
	// Get original invite
	invite, err := model.db.GetInvite(id)

	if err != nil {
		return err
	}

	// make sure user has access

	admins, err := model.db.GetOrgAdmins(invite.OrgId)

	if err != nil {
		return nil
	}

	isAdmin := false

	for _, admin := range admins {
		if admin.Id == userId {
			isAdmin = true
			break
		}
	}

	if isAdmin == false {
		return errors.New("Must be org admin to delete invite")
	}

	return model.db.DeleteInvite(id)
}
