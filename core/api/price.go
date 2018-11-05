package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
	"strconv"
	"time"
)

/**
 * @api {get} /org/:orgId/prices Get prices nearest in time or by currency
 * @apiVersion 1.0.1
 * @apiName GetPrices
 * @apiGroup Price
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiParam {Number} nearestDate Milliseconds since epoch
 * @apiParam {String} currency Currency code
 *
 * @apiSuccess {String} id Id of the Price.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {String} currency Currency code.
 * @apiSuccess {Date} date Date of the Price.
 * @apiSuccess {Date} inserted Date when Price was posted.
 * @apiSuccess {Date} updated Date when Price was updated.
 * @apiSuccess {Number} price Price of currency measured in native Org currency.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "11111111111111111111111111111111",
 *         "orgId": "11111111111111111111111111111111",
 *         "currency": "EUR",
 *         "date": "2018-09-11T18:05:04.420Z",
 *         "inserted": "2018-09-11T18:05:04.420Z",
 *         "updated": "2018-09-11T18:05:04.420Z",
 *         "price": 1.16
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetPrices(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	// TODO how do we make date an optional parameter
	// instead of resorting to this hack?
	nearestDate := time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC)

	nearestDateParam := r.URL.Query().Get("nearestDate")
	currencyParam := r.URL.Query().Get("currency")

	// If currency was specified, get all prices for that currency
	if currencyParam != "" {
		prices, err := model.Instance.GetPricesByCurrency(orgId, currencyParam, user.Id)

		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteJson(prices)
		return
	}

	if nearestDateParam != "" {
		nearestDateParamNumeric, err := strconv.ParseInt(nearestDateParam, 10, 64)

		if err != nil {
			rest.Error(w, "invalid date", 400)
			return
		}
		nearestDate = time.Unix(0, nearestDateParamNumeric*1000000)
	}

	// Get prices nearest in time
	prices, err := model.Instance.GetPricesNearestInTime(orgId, nearestDate, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(prices)
}

/**
 * @api {post} /orgs/:orgId/prices Create a new Price
 * @apiVersion 1.0.1
 * @apiName PostPrice
 * @apiGroup Price
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiParam {String} id Id 32 character hex string.
 * @apiParam {String} orgId Id of the Org.
 * @apiParam {String} currency Currency code.
 * @apiParam {Date} date Date of the Price.
 * @apiParam {Number} price Price of currency measured in native Org currency.
 *
 * @apiSuccess {String} id Id of the Price.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {String} currency Currency code.
 * @apiSuccess {Date} date Date of the Price.
 * @apiSuccess {Date} inserted Date when Price was posted.
 * @apiSuccess {Date} updated Date when Price was updated.
 * @apiSuccess {Number} price Price of currency measured in native Org currency.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "orgId": "11111111111111111111111111111111",
 *       "currency": "EUR",
 *       "date": "2018-09-11T18:05:04.420Z",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "price": 1.16
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostPrice(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	price := types.Price{}

	err := r.DecodeJsonPayload(&price)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	price.OrgId = orgId
	err = model.Instance.CreatePrice(&price, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&price)
}

/**
 * @api {delete} /orgs/:orgId/prices/:priceId Delete a Price
 * @apiVersion 1.0.1
 * @apiName DeletePrice
 * @apiGroup Price
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func DeletePrice(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	priceId := r.PathParam("priceId")

	err := model.Instance.DeletePrice(priceId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
