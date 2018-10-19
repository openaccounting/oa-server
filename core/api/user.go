package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/types"
	"net/http"
)

type VerifyUserParams struct {
	Code string `json:"code"`
}

type ConfirmResetPasswordParams struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type ResetPasswordParams struct {
	Email string `json:"email"`
}

/**
 * @api {get} /user Get Authenticated User
 * @apiVersion 1.0.0
 * @apiName GetUser
 * @apiGroup User
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiSuccess {String} id Id of the User.
 * @apiSuccess {Date} inserted Date User was created
 * @apiSuccess {Date} updated Date User was updated
 * @apiSuccess {String} firstName First name of the User.
 * @apiSuccess {String} lastName  Last name of the User.
 * @apiSuccess {String} email Email of the User.
 * @apiSuccess {Boolean} agreeToTerms Agree to terms
 * @apiSuccess {Boolean} emailVerified True if email has been verified.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "firstName": "John",
 *       "lastName": "Doe",
 *       "email": "johndoe@email.com",
 *       "agreeToTerms": true,
 *       "emailVerified": true
 *     }
 *
 * @apiUse NotAuthorizedError
 * @apiUse InternalServerError
 */
func GetUser(w rest.ResponseWriter, r *rest.Request) {
	user := r.Env["USER"].(*types.User)

	w.WriteJson(&user)
}

/**
 * @api {post} /users Create a new User
 * @apiVersion 1.0.0
 * @apiName PostUser
 * @apiGroup User
 *
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} id 32 character hex string
 * @apiParam {String} firstName First name of the User.
 * @apiParam {String} lastName  Last name of the User.
 * @apiParam {String} email Email of the User.
 * @apiParam {String} password Password of the User.
 * @apiParam {Boolean} agreeToTerms True if you agree to terms
 *
 * @apiSuccess {String} id Id of the User.
 * @apiSuccess {Date} inserted Date User was created
 * @apiSuccess {Date} updated Date User was updated
 * @apiSuccess {String} firstName First name of the User.
 * @apiSuccess {String} lastName  Last name of the User.
 * @apiSuccess {String} email Email of the User.
 * @apiSuccess {Boolean} emailVerified True if email has been verified.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "firstName": "John",
 *       "lastName": "Doe",
 *       "email": "johndoe@email.com",
 *       "agreeToTerms": true,
 *       "emailVerified": true
 *     }
 *
 * @apiUse InternalServerError
 */
func PostUser(w rest.ResponseWriter, r *rest.Request) {
	user := &types.User{}
	err := r.DecodeJsonPayload(user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.Instance.CreateUser(user)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(user)
}

/**
 * @api {put} /user Modify User
 * @apiVersion 1.0.0
 * @apiName PutUser
 * @apiGroup User
 *
 * @apiHeader {String} Authorization HTTP Basic Auth
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} password New password
 * @apiParam {String} code Password reset code. (Instead of Authorization header)
 *
 * @apiSuccess {String} id Id of the User.
 * @apiSuccess {Date} inserted Date User was created
 * @apiSuccess {Date} updated Date User was updated
 * @apiSuccess {String} firstName First name of the User.
 * @apiSuccess {String} lastName  Last name of the User.
 * @apiSuccess {String} email Email of the User.
 * @apiSuccess {Boolean} emailVerified True if email has been verified.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "id": "11111111111111111111111111111111",
 *       "inserted": "2018-09-11T18:05:04.420Z",
 *       "updated": "2018-09-11T18:05:04.420Z",
 *       "firstName": "John",
 *       "lastName": "Doe",
 *       "email": "johndoe@email.com",
 *       "agreeToTerms": true,
 *       "emailVerified": true
 *     }
 *
 * @apiUse InternalServerError
 */
func PutUser(w rest.ResponseWriter, r *rest.Request) {
	if r.Env["USER"] == nil {
		// password reset
		params := &ConfirmResetPasswordParams{}
		err := r.DecodeJsonPayload(params)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := model.Instance.ConfirmResetPassword(params.Password, params.Code)

		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteJson(user)
		return
	}

	// Otherwise it's an authenticated PUT

	user := r.Env["USER"].(*types.User)

	newUser := &types.User{}
	err := r.DecodeJsonPayload(newUser)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = newUser.Password

	err = model.Instance.UpdateUser(user)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(user)
}

/**
 * @api {post} /user/verify Verify user email address
 * @apiVersion 1.0.0
 * @apiName VerifyUser
 * @apiGroup User
 *
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} code Email verification code
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse InternalServerError
 */
func VerifyUser(w rest.ResponseWriter, r *rest.Request) {
	params := &VerifyUserParams{}

	err := r.DecodeJsonPayload(params)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.Instance.VerifyUser(params.Code)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/**
 * @api {post} /user/reset-password Send reset password email
 * @apiVersion 1.0.0
 * @apiName ResetPassword
 * @apiGroup User
 *
 * @apiHeader {String} Accept-Version ^1.0.0 semver versioning
 *
 * @apiParam {String} email Email address for user
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse InternalServerError
 */
func ResetPassword(w rest.ResponseWriter, r *rest.Request) {
	params := &ResetPasswordParams{}

	err := r.DecodeJsonPayload(params)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = model.Instance.ResetPassword(params.Email)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
