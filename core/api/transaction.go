package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
)

/**
 * @api {get} /orgs/:orgId/accounts/:accountId/transactions Get Transactions by Account Id
 * @apiVersion 1.0.1
 * @apiName GetAccountTransactions
 * @apiGroup Transaction
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiSuccess {String} id Id of the Transaction.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {String} userId Id of the User who created the Transaction.
 * @apiSuccess {Date} date Date of the Transaction
 * @apiSuccess {Date} inserted Date Transaction was created
 * @apiSuccess {Date} updated Date Transaction was updated
 * @apiSuccess {String} description Description of Transaction
 * @apiSuccess {String} data Extra data field
 * @apiSuccess {Object[]} splits Array of Transaction Splits
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "11111111111111111111111111111111",
 *         "orgId": "11111111111111111111111111111111",
 *         "userId": "11111111111111111111111111111111",
 *         "date": "2018-06-08T20:12:29.720Z",
 *         "inserted": "2018-06-08T20:12:29.720Z",
 *         "updated": "2018-06-08T20:12:29.720Z",
 *         "description": "Treat friend to lunch",
 *         "data:": "{\"key\": \"value\"}",
 *         "splits": [
 *           {
 *             "accountId": "11111111111111111111111111111111",
 *             "amount": -2000,
 *             "nativeAmount": -2000
 *           },
 *           {
 *             "accountId": "22222222222222222222222222222222",
 *             "amount": 1000,
 *             "nativeAmount": 1000
 *           },
 *           {
 *             "accountId": "33333333333333333333333333333333",
 *             "amount": 1000,
 *             "nativeAmount": 1000
 *           }
 *         ]
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetTransactionsByAccount(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")
	accountId := r.PathParam("accountId")

	queryOptions, err := types.QueryOptionsFromURLQuery(r.URL.Query())

	if err != nil {
		rest.Error(w, "invalid query options", 400)
		return
	}

	sTxs, err := model.Instance.GetTransactionsByAccount(orgId, user.Id, accountId, queryOptions)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&sTxs)
}

/**
 * @api {get} /orgs/:orgId/transactions Get Transactions by Org Id
 * @apiVersion 1.0.1
 * @apiName GetOrgTransactions
 * @apiGroup Transaction
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiSuccess {String} id Id of the Transaction.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {String} userId Id of the User who created the Transaction.
 * @apiSuccess {Date} date Date of the Transaction
 * @apiSuccess {Date} inserted Date Transaction was created
 * @apiSuccess {Date} updated Date Transaction was updated
 * @apiSuccess {String} description Description of Transaction
 * @apiSuccess {String} data Extra data field
 * @apiSuccess {Object[]} splits Array of Transaction Splits
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "11111111111111111111111111111111",
 *         "orgId": "11111111111111111111111111111111",
 *         "userId": "11111111111111111111111111111111",
 *         "date": "2018-06-08T20:12:29.720Z",
 *         "inserted": "2018-06-08T20:12:29.720Z",
 *         "updated": "2018-06-08T20:12:29.720Z",
 *         "description": "Treat friend to lunch",
 *         "data:": "{\"key\": \"value\"}",
 *         "splits": [
 *           {
 *             "accountId": "11111111111111111111111111111111",
 *             "amount": -2000,
 *             "nativeAmount": -2000
 *           },
 *           {
 *             "accountId": "22222222222222222222222222222222",
 *             "amount": 1000,
 *             "nativeAmount": 1000
 *           },
 *           {
 *             "accountId": "33333333333333333333333333333333",
 *             "amount": 1000,
 *             "nativeAmount": 1000
 *           }
 *         ]
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetTransactionsByOrg(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	queryOptions, err := types.QueryOptionsFromURLQuery(r.URL.Query())

	if err != nil {
		rest.Error(w, "invalid query options", 400)
		return
	}

	sTxs, err := model.Instance.GetTransactionsByOrg(orgId, user.Id, queryOptions)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&sTxs)
}

/**
 * @api {post} /orgs/:orgId/transactions Create a new Transaction
 * @apiVersion 1.0.1
 * @apiName PostTransaction
 * @apiGroup Transaction
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiParam {String} id Id 32 character hex string
 * @apiParam {Date} date Date of the Transaction
 * @apiParam {String} description Description of Transaction
 * @apiParam {String} data Extra data field
 * @apiParam {Object[]} splits Array of Transaction Splits. nativeAmounts must add up to 0.
 * @apiParam {String} splits.accountId Id of Account
 * @apiParam {Number} splits.amount Amount of split in Account currency
 * @apiParam {Number} splits.nativeAmount Amount of split in Org currency
 *
 * @apiSuccess {String} id Id of the Transaction.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {String} userId Id of the User who created the Transaction.
 * @apiSuccess {Date} date Date of the Transaction
 * @apiSuccess {Date} inserted Date Transaction was created
 * @apiSuccess {Date} updated Date Transaction was updated
 * @apiSuccess {String} description Description of Transaction
 * @apiSuccess {String} data Extra data field
 * @apiSuccess {Object[]} splits Array of Transaction Splits
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "orgId": "11111111111111111111111111111111",
 *       "userId": "11111111111111111111111111111111",
 *       "date": "2018-06-08T20:12:29.720Z",
 *       "inserted": "2018-06-08T20:12:29.720Z",
 *       "updated": "2018-06-08T20:12:29.720Z",
 *       "description": "Treat friend to lunch",
 *       "data:": "{\"key\": \"value\"}",
 *       "splits": [
 *         {
 *           "accountId": "11111111111111111111111111111111",
 *           "amount": -2000,
 *           "nativeAmount": -2000
 *         },
 *         {
 *           "accountId": "22222222222222222222222222222222",
 *           "amount": 1000,
 *           "nativeAmount": 1000
 *         },
 *         {
 *           "accountId": "33333333333333333333333333333333",
 *           "amount": 1000,
 *           "nativeAmount": 1000
 *         }
 *       ]
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostTransaction(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	sTx := types.Transaction{}
	err := r.DecodeJsonPayload(&sTx)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sTx.OrgId = orgId
	sTx.UserId = user.Id

	err = model.Instance.CreateTransaction(&sTx)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(sTx)
}

/**
 * @api {put} /orgs/:orgId/transactions/:transactionId Modify a Transaction
 * @apiVersion 1.0.1
 * @apiName PutTransaction
 * @apiGroup Transaction
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiParam {String} id 32 character hex string
 * @apiParam {Date} date Date of the Transaction
 * @apiParam {String} description Description of Transaction
 * @apiParam {String} data Extra data field
 * @apiParam {Object[]} splits Array of Transaction Splits. nativeAmounts must add up to 0.
 * @apiParam {String} splits.accountId Id of Account
 * @apiParam {Number} splits.amount Amount of split in Account currency
 * @apiParam {Number} splits.nativeAmount Amount of split in Org currency
 *
 * @apiSuccess {String} id Id of the Transaction.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {String} userId Id of the User who created the Transaction.
 * @apiSuccess {Date} date Date of the Transaction
 * @apiSuccess {Date} inserted Date Transaction was created
 * @apiSuccess {Date} updated Date Transaction was updated
 * @apiSuccess {String} description Description of Transaction
 * @apiSuccess {String} data Extra data field
 * @apiSuccess {Object[]} splits Array of Transaction Splits
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "orgId": "11111111111111111111111111111111",
 *       "userId": "11111111111111111111111111111111",
 *       "date": "2018-06-08T20:12:29.720Z",
 *       "inserted": "2018-06-08T20:12:29.720Z",
 *       "updated": "2018-06-08T20:12:29.720Z",
 *       "description": "Treat friend to lunch",
 *       "data:": "{\"key\": \"value\"}",
 *       "splits": [
 *         {
 *           "accountId": "11111111111111111111111111111111",
 *           "amount": -2000,
 *           "nativeAmount": -2000
 *         },
 *         {
 *           "accountId": "22222222222222222222222222222222",
 *           "amount": 1000,
 *           "nativeAmount": 1000
 *         },
 *         {
 *           "accountId": "33333333333333333333333333333333",
 *           "amount": 1000,
 *           "nativeAmount": 1000
 *         }
 *       ]
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PutTransaction(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")
	transactionId := r.PathParam("transactionId")

	sTx := types.Transaction{}
	err := r.DecodeJsonPayload(&sTx)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sTx.OrgId = orgId
	sTx.UserId = user.Id

	err = model.Instance.UpdateTransaction(transactionId, &sTx)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(sTx)
}

/**
 * @api {delete} /orgs/:orgId/transactions/:transactionId Delete a Transaction
 * @apiVersion 1.0.1
 * @apiName DeleteTransaction
 * @apiGroup Transaction
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
func DeleteTransaction(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")
	transactionId := r.PathParam("transactionId")

	err := model.Instance.DeleteTransaction(transactionId, user.Id, orgId)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
