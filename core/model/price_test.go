package model

import (
	"errors"
	"github.com/openaccounting/oa-server/core/mocks"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreatePrice(t *testing.T) {

	price := types.Price{
		"1",
		"2",
		"BTC",
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Unix(0, 0),
		6700,
	}

	badPrice := types.Price{
		"1",
		"2",
		"",
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Unix(0, 0),
		6700,
	}

	badOrg := types.Price{
		"1",
		"1",
		"BTC",
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Unix(0, 0),
		6700,
	}

	tests := map[string]struct {
		err   error
		price types.Price
	}{
		"successful": {
			err:   nil,
			price: price,
		},
		"with error": {
			err:   errors.New("currency required"),
			price: badPrice,
		},
		"with org error": {
			err:   errors.New("User does not belong to org"),
			price: badOrg,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		price := test.price
		userId := "3"

		db := &mocks.Datastore{}

		db.On("GetOrgs", userId).Return([]*types.Org{
			{
				Id: "2",
			},
		}, nil)

		db.On("InsertPrice", &test.price).Return(nil)

		db.On("GetOrgUserIds", price.OrgId).Return([]string{userId}, nil)

		model := NewModel(db, &util.StandardBcrypt{}, types.Config{})

		err := model.CreatePrice(&price, userId)

		assert.Equal(t, test.err, err)
	}
}

func TestDeletePrice(t *testing.T) {

	price := types.Price{
		"1",
		"2",
		"BTC",
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Unix(0, 0),
		6700,
	}

	tests := map[string]struct {
		err    error
		userId string
		price  types.Price
	}{
		"successful": {
			err:    nil,
			price:  price,
			userId: "3",
		},
		"with org error": {
			err:    errors.New("User does not belong to org"),
			price:  price,
			userId: "4",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		price := test.price

		db := &mocks.Datastore{}

		db.On("GetPriceById", price.Id).Return(&price, nil)

		db.On("GetOrgs", "3").Return([]*types.Org{
			{
				Id: "2",
			},
		}, nil)

		db.On("GetOrgs", "4").Return([]*types.Org{
			{
				Id: "7",
			},
		}, nil)

		db.On("DeletePrice", price.Id).Return(nil)

		db.On("GetOrgUserIds", price.OrgId).Return([]string{test.userId}, nil)

		model := NewModel(db, &util.StandardBcrypt{}, types.Config{})

		err := model.DeletePrice(price.Id, test.userId)

		assert.Equal(t, test.err, err)
	}
}
