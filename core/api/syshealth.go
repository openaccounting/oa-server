package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model"
)

/**
 * @api {get} /health-check Get system health status
 * @apiVersion 1.2.0
 * @apiName GetSystemHealthStatus
 * @apiGroup SystemHealth
 *
 *
 * @apiHeader {String} Accept-Version: 1.2.0 semver versioning
 *
 * @apiSuccess {String} database Database status: "ok"; "fail"
 * @apiSuccess {String} api API status: "ok"
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "database": "ok",
 *       "api": "ok",
 *     }
 *
 * @apiUse InternalServerError
 */
func GetSystemHealthStatus(w rest.ResponseWriter, r *rest.Request) {
	status := map[string]string{
		"database": "ok",
		"api":      "ok",
	}
	if err := model.Instance.PingDatabase(); err != nil {
		status["database"] = "fail"
	}
	w.WriteJson(status)
}
