package api

import (
	"github.com/Masterminds/semver"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

type VersionMiddleware struct {
}

// MiddlewareFunc makes AuthMiddleware implement the Middleware interface.
func (mw *VersionMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {
		version := request.Header.Get("Accept-Version")

		// Don't require version header for websockets
		if request.URL.String() == "/ws" {
			handler(writer, request)
			return
		}

		if version == "" {
			rest.Error(writer, "Accept-Version header required", http.StatusBadRequest)
			return
		}

		constraint, err := semver.NewConstraint(version)

		if err != nil {
			rest.Error(writer, "Invalid version", http.StatusBadRequest)
		}

		serverVersion, _ := semver.NewVersion("1.3.0")
		// Pre-release versions
		compatVersion, _ := semver.NewVersion("0.1.8")

		versionMatch := constraint.Check(serverVersion)
		compatMatch := constraint.Check(compatVersion)

		if versionMatch == false && compatMatch == false {
			rest.Error(writer, "Invalid version", http.StatusBadRequest)
			return
		}

		handler(writer, request)
	}
}
