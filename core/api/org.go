package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
)

/**
 * @api {get} /org/:orgId Get Org by id
 * @apiVersion 1.0.0
 * @apiName GetOrg
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiSuccess {String} id Id of the Org.
 * @apiSuccess {Date} inserted Date Org was created
 * @apiSuccess {Date} updated Date Org was updated
 * @apiSuccess {String} name Name of the Org.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "name": "MyOrg",
 *       "currency": "USD",
 *       "precision": 2,
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetOrg(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	org, err := model.Instance.GetOrg(orgId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&org)
}

/**
 * @api {get} /orgs Get a User's Orgs
 * @apiVersion 1.0.0
 * @apiName GetOrgs
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiSuccess {String} id Id of the Org.
 * @apiSuccess {Date} inserted Date Org was created
 * @apiSuccess {Date} updated Date Org was updated
 * @apiSuccess {String} name Name of the Org.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "11111111111111111111111111111111",
 *         "inserted": "2018-09-11T18:05:04.420Z",
 *         "updated": "2018-09-11T18:05:04.420Z",
 *         "name": "MyOrg",
 *         "currency": "USD",
 *         "precision": 2,
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetOrgs(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)

	orgs, err := model.Instance.GetOrgs(user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&orgs)
}

/**
 * @api {post} /orgs Create a new Org
 * @apiVersion 1.0.0
 * @apiName PostOrg
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} id Id 32 character hex string
 * @apiParam {String} name Name of the Org.
 * @apiParam {String} currency Three letter currency code.
 * @apiParam {Number} precision How many digits the currency goes out to.
 *
 * @apiSuccess {String} id Id of the Org.
 * @apiSuccess {Date} inserted Date Org was created
 * @apiSuccess {Date} updated Date Org was updated
 * @apiSuccess {String} name Name of the Org.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "name": "MyOrg",
 *       "currency": "USD",
 *       "precision": 2,
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostOrg(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	org := types.Org{Precision: 2}
	err := r.DecodeJsonPayload(&org)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.Instance.CreateOrg(&org, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&org)
}

/**
 * @api {put} /orgs/:orgId Modify an Org
 * @apiVersion 1.0.0
 * @apiName PutOrg
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} name Name of the Org.
 *
 * @apiSuccess {String} id Id of the Org.
 * @apiSuccess {Date} inserted Date Org was created
 * @apiSuccess {Date} updated Date Org was updated
 * @apiSuccess {String} name Name of the Org.
 * @apiSuccess {String} currency Three letter currency code.
 * @apiSuccess {Number} precision How many digits the currency goes out to.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "name": "MyOrg",
 *       "currency": "USD",
 *       "precision": 2,
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PutOrg(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	org := types.Org{}
	err := r.DecodeJsonPayload(&org)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	org.Id = orgId

	err = model.Instance.UpdateOrg(&org, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&org)
}

/**
 * @api {post} /orgs/:orgId/invites Invite a user to an Org
 * @apiVersion 1.0.0
 * @apiName PostInvite
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} email Email address of user
 *
 * @apiSuccess {String} id Id of the Invite
 * @apiSuccess {orgId} id Id of the Org
 * @apiSuccess {Date} inserted Date Invite was created
 * @apiSuccess {Date} updated Date Invite was updated/accepted
 * @apiSuccess {String} email Email address of user
 * @apiSuccess {String} accepted true if user has accepted
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "a1b2c3d4",
 *       "orgId": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "email": "johndoe@email.com",
 *       "accepted": false
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PostInvite(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	invite := types.Invite{}
	err := r.DecodeJsonPayload(&invite)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	invite.OrgId = orgId

	err = model.Instance.CreateInvite(&invite, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&invite)
}

/**
 * @api {put} /orgs/:orgId/invites/:inviteId Accept an invitation
 * @apiVersion 1.0.0
 * @apiName PutInvite
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} accepted true
 *
 * @apiSuccess {String} id Id of the Invite
 * @apiSuccess {orgId} id Id of the Org
 * @apiSuccess {Date} inserted Date Invite was created
 * @apiSuccess {Date} updated Date Invite was updated/accepted
 * @apiSuccess {String} email Email address of user
 * @apiSuccess {String} accepted true if user has accepted
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "a1b2c3d4",
 *       "orgId": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "email": "johndoe@email.com",
 *       "accepted": true
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func PutInvite(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	//orgId := r.PathParam("orgId")
	inviteId := r.PathParam("inviteId")

	invite := types.Invite{}
	err := r.DecodeJsonPayload(&invite)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	invite.Id = inviteId

	err = model.Instance.AcceptInvite(&invite, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&invite)
}

/**
 * @api {get} /orgs/:orgId/invites Get Org invites
 * @apiVersion 1.0.0
 * @apiName GetInvites
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiSuccess {String} id Id of the Invite
 * @apiSuccess {orgId} id Id of the Org
 * @apiSuccess {Date} inserted Date Invite was created
 * @apiSuccess {Date} updated Date Invite was updated/accepted
 * @apiSuccess {String} email Email address of user
 * @apiSuccess {String} accepted true if user has accepted
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     [
 *       {
 *         "id": "a1b2c3d4",
 *         "orgId": "11111111111111111111111111111111",
 *         "inserted": "2018-09-11T18:05:04.420Z",
 *         "updated": "2018-09-11T18:05:04.420Z",
 *         "email": "johndoe@email.com",
 *         "accepted": true
 *       }
 *     ]
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetInvites(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	orgId := r.PathParam("orgId")

	invites, err := model.Instance.GetInvites(orgId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&invites)
}

/**
 * @api {delete} /orgs/:orgId/invites/:inviteId Delete Invite
 * @apiVersion 1.0.0
 * @apiName DeleteInvite
 * @apiGroup Org
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func DeleteInvite(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)
	inviteId := r.PathParam("inviteId")

	err := model.Instance.DeleteInvite(inviteId, user.Id)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
