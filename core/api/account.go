package api

import (
	"encoding/json"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

/**
 * @api {get} /orgs/:orgId/accounts Get Accounts by Org id
 * @apiVersion 1.0.1
 * @apiName GetOrgAccounts
 * @apiGroup Account
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiSuccess {String} id Id of the Account.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {Date} inserted Date Account was created
 * @apiSuccess {Date} updated Date Account was updated
 * @apiSuccess {String} name Name of the Account.
 * @apiSuccess {String} parent Id of the parent Account.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 * @apiSuccess {Boolean} debitBalance True if Account has a debit balance.
 * @apiSuccess {Number} balance Current Account balance in this Account's currency
 * @apiSuccess {Number} nativeBalance Current Account balance in the Org's currency
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "22222222222222222222222222222222",
 *         "orgId": "11111111111111111111111111111111",
 *         "inserted": "2018-09-11T18:05:04.420Z",
 *         "updated": "2018-09-11T18:05:04.420Z",
 *         "name": "Cash",
 *         "parent": "11111111111111111111111111111111",
 *         "currency": "USD",
 *         "precision": 2,
 *         "debitBalance": true,
 *         "balance": 10000,
 *         "nativeBalance": 10000
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetOrgAccounts(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	// TODO how do we make date an optional parameter
	// instead of resorting to this hack?
	date := time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC)

	dateParam := r.URL.Query().Get("date")

	if dateParam != "" {
		dateParamNumeric, err := strconv.ParseInt(dateParam, 10, 64)

		if err != nil {
			rest.Error(w, "invalid date", 400)
			return
		}
		date = time.Unix(0, dateParamNumeric*1000000)
	}

	accounts, err := model.Instance.GetAccountsWithBalances(orgId, user.Id, "", date)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&accounts)
}

/**
 * @api {post} /orgs/:orgId/accounts Create a new Account
 * @apiVersion 1.0.1
 * @apiName PostAccount
 * @apiGroup Account
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiParam {String} id Id 32 character hex string
 * @apiParam {String} name Name of the Account.
 * @apiParam {String} parent Id of the parent Account.
 * @apiParam {String} currency Three letter currency code.
 * @apiParam {Number} precision How many digits the currency goes out to.
 * @apiParam {Boolean} debitBalance True if account has a debit balance.
 * @apiParam {Number} balance Current Account balance in this Account's currency
 * @apiParam {Number} nativeBalance Current Account balance in the Org's currency
 *
 * @apiSuccess {String} id Id of the Account.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {Date} inserted Date Account was created
 * @apiSuccess {Date} updated Date Account was updated
 * @apiSuccess {String} name Name of the Account.
 * @apiSuccess {String} parent Id of the parent Account.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 * @apiSuccess {Boolean} debitBalance True if account has a debit balance.
 * @apiSuccess {Number} balance Current Account balance in this Account's currency
 * @apiSuccess {Number} nativeBalance Current Account balance in the Org's currency
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "22222222222222222222222222222222",
 *       "orgId": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "name": "Cash",
 *       "parent": "11111111111111111111111111111111",
 *       "currency": "USD",
 *       "precision": 2,
 *       "debitBalance": true,
 *       "balance": 10000,
 *       "nativeBalance": 10000
 *       }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostAccount(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(content) == 0 {
		rest.Error(w, "JSON payload is empty", http.StatusInternalServerError)
		return
	}

	account := types.NewAccount()

	err = json.Unmarshal(content, &account)

	if err != nil {
		// Maybe it's an array of accounts?
		PostAccounts(w, r, content)
		return
	}

	account.OrgId = orgId
	err = model.Instance.CreateAccount(account, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&account)
}

func PostAccounts(w rest.ResponseWriter, r *rest.Request, content []byte) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	accounts := make([]*types.Account, 0)

	err := json.Unmarshal(content, &accounts)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, account := range accounts {
		account.OrgId = orgId
		err = model.Instance.CreateAccount(account, user.Id)

		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteJson(accounts)
}

/**
 * @api {put} /orgs/:orgId/accounts/:accountId Modify an Account
 * @apiVersion 1.0.1
 * @apiName PutAccount
 * @apiGroup Account
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 *
 * @apiParam {String} id Id 32 character hex string
 * @apiParam {String} name Name of the Account.
 * @apiParam {String} parent Id of the parent Account.
 * @apiParam {String} currency Three letter currency code.
 * @apiParam {Number} precision How many digits the currency goes out to.
 * @apiParam {Boolean} debitBalance True if Account has a debit balance.
 * @apiParam {Number} balance Current Account balance in this Account's currency
 * @apiParam {Number} nativeBalance Current Account balance in the Org's currency
 *
 * @apiSuccess {String} id Id of the Account.
 * @apiSuccess {String} orgId Id of the Org.
 * @apiSuccess {Date} inserted Date Account was created
 * @apiSuccess {Date} updated Date Account was updated
 * @apiSuccess {String} name Name of the Account.
 * @apiSuccess {String} parent Id of the parent Account.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 * @apiSuccess {Boolean} debitBalance True if Account has a debit balance.
 * @apiSuccess {Number} balance Current Account balance in this Account's currency
 * @apiSuccess {Number} nativeBalance Current Account balance in the Org's currency
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "22222222222222222222222222222222",
 *       "orgId": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *        "updated": "2018-09-11T18:05:04.420Z",
 *       "name": "Cash",
 *       "parent": "11111111111111111111111111111111",
 *       "currency": "USD",
 *       "precision": 2,
 *       "debitBalance": true,
 *       "balance": 10000,
 *       "nativeBalance": 10000
 *       }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PutAccount(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")
	accountId := r.PathParam("accountId")

	account := types.Account{}
	err := r.DecodeJsonPayload(&account)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account.Id = accountId
	account.OrgId = orgId

	err = model.Instance.UpdateAccount(&account, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&account)
}

/**
 * @api {delete} /orgs/:orgId/accounts/:accountId Delete an Account
 * @apiVersion 1.0.1
 * @apiName DeleteAccount
 * @apiGroup Account
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
func DeleteAccount(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")
	accountId := r.PathParam("accountId")

	err := model.Instance.DeleteAccount(accountId, user.Id, orgId)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
