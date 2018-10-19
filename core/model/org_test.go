package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TdOrg struct {
	db.Datastore
}

func (td *TdOrg) GetOrg(orgId string, userId string) (*types.Org, error) {
	if userId == "1" {
		return &types.Org{
			Id:        "1",
			Name:      "MyOrg",
			Currency:  "USD",
			Precision: 2,
		}, nil
	} else {
		return nil, errors.New("not found")
	}
}

func (td *TdOrg) UpdateOrg(org *types.Org) error {
	return nil
}

func TestUpdateOrg(t *testing.T) {
	tests := map[string]struct {
		err    error
		org    *types.Org
		userId string
	}{
		"success": {
			err: nil,
			org: &types.Org{
				Id:   "1",
				Name: "MyOrg2",
			},
			userId: "1",
		},
		"access denied": {
			err: errors.New("access denied"),
			org: &types.Org{
				Id:   "1",
				Name: "MyOrg2",
			},
			userId: "2",
		},
		"error": {
			err: errors.New("name required"),
			org: &types.Org{
				Id:   "1",
				Name: "",
			},
			userId: "1",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		td := &TdOrg{}

		model := NewModel(td, nil, types.Config{})

		err := model.UpdateOrg(test.org, test.userId)
		assert.Equal(t, test.err, err)
	}
}
