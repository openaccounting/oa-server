package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
)

/**
 * @api {get} /orgs/:orgId/budget Get Budget
 * @apiVersion 1.4.0
 * @apiName GetBudget
 * @apiGroup Budget
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.4.0 semver versioning
 *
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {Date} inserted Date Transaction was created
 * @apiSuccess {Object[]} items Array of Budget Items
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "orgId": "11111111111111111111111111111111",
 *         "inserted": "2020-01-13T20:12:29.720Z",
 *         "items": [
 *           {
 *             "accountId": "11111111111111111111111111111111",
 *             "amount": 35000,
 *           },
 *           {
 *             "accountId": "22222222222222222222222222222222",
 *             "amount": 55000
 *           }
 *         ]
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetBudget(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	budget, err := model.Instance.GetBudget(orgId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&budget)
}

/**
 * @api {post} /orgs/:orgId/budget Create a Budget
 * @apiVersion 1.4.0
 * @apiName PostBudget
 * @apiGroup Budget
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.4.0 semver versioning
 *
 * @apiParam {Object[]} items Array of Budget Items.
 * @apiParam {String} items.accountId Id of Expense Account
 * @apiParam {Number} items.amount Amount budgeted
 *
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {Date} inserted Date Transaction was created
 * @apiSuccess {Object[]} items Array of Budget Items
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "orgId": "11111111111111111111111111111111",
 *         "inserted": "2020-01-13T20:12:29.720Z",
 *         "items": [
 *           {
 *             "accountId": "11111111111111111111111111111111",
 *             "amount": 35000,
 *           },
 *           {
 *             "accountId": "22222222222222222222222222222222",
 *             "amount": 55000
 *           }
 *         ]
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostBudget(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	budget := types.Budget{}
	err := r.DecodeJsonPayload(&budget)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	budget.OrgId = orgId

	err = model.Instance.CreateBudget(&budget, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(budget)
}

/**
 * @api {delete} /orgs/:orgId/budget Delete Budget
 * @apiVersion 1.4.0
 * @apiName DeleteBudget
 * @apiGroup Budget
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.4.0 semver versioning
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func DeleteBudget(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	err := model.Instance.DeleteBudget(orgId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
