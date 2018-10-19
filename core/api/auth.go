package api

import (
	"encoding/base64"
	"errors"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/auth"
	"log"
	"net/http"
	"strings"
)

type AuthMiddleware struct {

	// Realm name to display to the user. Required.
	Realm string
}

// MiddlewareFunc makes AuthMiddleware implement the Middleware interface.
func (mw *AuthMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {

	if mw.Realm == "" {
		log.Fatal("Realm is required")
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {

		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			request.Env["USER"] = nil
			handler(writer, request)
			return
		}

		emailOrKey, password, err := mw.decodeBasicAuthHeader(authHeader)

		if err != nil {
			rest.Error(writer, "Invalid authentication", http.StatusBadRequest)
			return
		}

		// authenticate via session, apikey or user
		user, err := auth.Instance.Authenticate(emailOrKey, password)

		if err == nil {
			request.Env["USER"] = user
			handler(writer, request)
			return
		}

		log.Println("Unauthorized " + emailOrKey)

		mw.unauthorized(writer)
		return
	}
}

func (mw *AuthMiddleware) unauthorized(writer rest.ResponseWriter) {
	writer.Header().Set("WWW-Authenticate", "Basic realm="+mw.Realm)
	rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}

func (mw *AuthMiddleware) decodeBasicAuthHeader(header string) (user string, password string, err error) {

	parts := strings.SplitN(header, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Basic") {
		return "", "", errors.New("Invalid authentication")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("Invalid base64")
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", errors.New("Invalid authentication")
	}

	return creds[0], creds[1], nil
}

func (mw *AuthMiddleware) RequireAuth(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(writer rest.ResponseWriter, request *rest.Request) {

		if request.Env["USER"] == nil {
			mw.unauthorized(writer)
			return
		}

		handler(writer, request)
	}
}
