package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
)

/**
 * @api {post} /sessions Create a new Session
 * @apiVersion 1.0.1
 * @apiName PostSession
 * @apiGroup Session
 *
 * @apiHeader {String} Accept-Version ^1.0.1 semver versioning
 * @apiHeader {String} Authorization HTTP Basic Auth
 *
 * @apiParam {String} id 32 character hex string
 *
 * @apiSuccess {String} id Id of the Session.
 * @apiSuccess {Date} inserted Date Session was created
 * @apiSuccess {Date} updated Date Last activity for the Session
 * @apiSuccess {String} userId Id of the User
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "userId": "22222222222222222222222222222222"
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostSession(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	session := &types.Session{}

	err := r.DecodeJsonPayload(session)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.UserId = user.Id

	err = model.Instance.CreateSession(session)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(session)
}

/**
 * @api {delete} /sessions/:sessionId Log out of a Session
 * @apiVersion 1.0.1
 * @apiName DeleteSession
 * @apiGroup Session
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
func DeleteSession(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	sessionId := r.PathParam("sessionId")

	err := model.Instance.DeleteSession(sessionId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
