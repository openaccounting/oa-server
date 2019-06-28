package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
)

/**
 * @api {get} /apikeys Get API keys
 * @apiVersion 1.3.0
 * @apiName GetApiKeys
 * @apiGroup ApiKey
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.3.0 semver versioning
 *
 * @apiSuccess {String} id Id of the ApiKey.
 * @apiSuccess {Date} inserted Date ApiKey was created
 * @apiSuccess {Date} updated Date Last activity for the ApiKey
 * @apiSuccess {String} userId Id of the User
 * @apiSuccess {String} label Label
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "11111111111111111111111111111111",
 *         "inserted": "2018-09-11T18:05:04.420Z",
 *         "updated": "2018-09-11T18:05:04.420Z",
 *         "userId": "22222222222222222222222222222222",
 *         "label": "Shopping Cart"
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetApiKeys(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)

	keys, err := model.Instance.GetApiKeys(user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(keys)
}

/**
 * @api {post} /apikeys Create a new API key
 * @apiVersion 1.3.0
 * @apiName PostApiKey
 * @apiGroup ApiKey
 *
 * @apiHeader {String} Accept-Version ^1.3.0 semver versioning
 * @apiHeader {String} Authorization HTTP Basic Auth
 *
 * @apiParam {String} id 32 character hex string
 * @apiParam {String} label Label
 *
 * @apiSuccess {String} id Id of the ApiKey.
 * @apiSuccess {Date} inserted Date ApiKey was created
 * @apiSuccess {Date} updated Date Last activity for the ApiKey
 * @apiSuccess {String} userId Id of the User
 * @apiSuccess {String} label Label
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "userId": "22222222222222222222222222222222",
 *       "label": "Shopping Cart"
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostApiKey(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	key := &types.ApiKey{}

	err := r.DecodeJsonPayload(key)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key.UserId = user.Id

	err = model.Instance.CreateApiKey(key)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(key)
}

/**
 * @api {put} /apikeys Modify an API key
 * @apiVersion 1.3.0
 * @apiName PutApiKey
 * @apiGroup ApiKey
 *
 * @apiHeader {String} Accept-Version ^1.3.0 semver versioning
 * @apiHeader {String} Authorization HTTP Basic Auth
 *
 * @apiParam {String} id 32 character hex string
 * @apiParam {String} label Label
 *
 * @apiSuccess {String} id Id of the ApiKey.
 * @apiSuccess {Date} inserted Date ApiKey was created
 * @apiSuccess {Date} updated Date Last activity for the ApiKey
 * @apiSuccess {String} userId Id of the User
 * @apiSuccess {String} label Label
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "userId": "22222222222222222222222222222222",
 *       "label": "Shopping Cart"
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PutApiKey(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	key := &types.ApiKey{}
	keyId := r.PathParam("apiKeyId")

	err := r.DecodeJsonPayload(key)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key.Id = keyId
	key.UserId = user.Id

	err = model.Instance.UpdateApiKey(key)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(key)
}

/**
 * @api {delete} /apikeys/:apiKeyId Delete an API key
 * @apiVersion 1.3.0
 * @apiName DeleteApiKey
 * @apiGroup ApiKey
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.3.0 semver versioning
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func DeleteApiKey(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	id := r.PathParam("apiKeyId")

	err := model.Instance.DeleteApiKey(id, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
