package api

import (
	"github.com/ant0ine/go-json-rest/rest"
)

/**
 * Changelog
 *
 * 1.0.1
 * - add user.signupSource
 *
 */

/**
 * @apiDefine NotAuthorizedError
 *
 * @apiError NotAuthorized API request does not have proper credentials
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 403 Not Authorized
 */

/**
 * @apiDefine InternalServerError
 *
 * @apiError InternalServer An internal error occurred
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 500 Internal Server Error
 *     {
 *       "error": "id required"
 *     }
 *
 */

func Init(prefix string) (*rest.Api, error) {
	rest.ErrorFieldName = "error"
	app := rest.NewApi()

	logger := &LoggerMiddleware{}

	var stack = []rest.Middleware{
		logger,
		&rest.RecorderMiddleware{},
		&rest.TimerMiddleware{},
		&rest.PoweredByMiddleware{},
		&rest.RecoverMiddleware{},
		&rest.GzipMiddleware{},
		&rest.ContentTypeCheckerMiddleware{},
	}

	app.Use(stack...)

	app.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			//return origin == "http://localhost:4200"
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin", "Authorization", "Accept-Version"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	auth := &AuthMiddleware{
		Realm: "openaccounting",
	}

	version := &VersionMiddleware{}

	app.Use(auth)
	app.Use(version)

	router, err := GetRouter(auth, prefix)
	if err != nil {
		return nil, err
	}

	app.SetApp(router)

	return app, nil
}
