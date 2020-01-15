package api

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/openaccounting/oa-server/core/model/types"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type LoggerMiddleware struct {
	Logger *log.Logger
}

func (mw *LoggerMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {

	// set the default Logger
	if mw.Logger == nil {
		mw.Logger = log.New(os.Stderr, "", 0)
	}

	return func(w rest.ResponseWriter, r *rest.Request) {
		h(w, r)

		message := getIp(r)

		message = message + " " + getUser(r)
		message = message + " " + getTime(r)
		message = message + " " + getRequest(r)
		message = message + " " + getStatus(r)
		message = message + " " + getBytes(r)
		message = message + " " + getUserAgent(r)
		message = message + " " + getElapsedTime(r)

		mw.Logger.Print(message)
	}
}

func getIp(r *rest.Request) string {
	remoteAddr := r.RemoteAddr
	if remoteAddr != "" {
		if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
			return ip
		}
	}
	return ""
}

func getUser(r *rest.Request) string {
	if r.Env["USER"] != nil {
		user := r.Env["USER"].(*types.User)
		return user.Email
	}

	return "-"
}

func getTime(r *rest.Request) string {
	if r.Env["START_TIME"] != nil {
		return r.Env["START_TIME"].(*time.Time).Format("02/Jan/2006:15:04:05 -0700")
	}
	return "-"
}

func getElapsedTime(r *rest.Request) string {
	if r.Env["ELAPSED_TIME"] != nil {
		return r.Env["ELAPSED_TIME"].(*time.Duration).String()
	}
	return "-"
}

func getRequest(r *rest.Request) string {
	return r.Method + " " + r.URL.RequestURI()
}

func getStatus(r *rest.Request) string {
	if r.Env["STATUS_CODE"] != nil {
		return strconv.Itoa(r.Env["STATUS_CODE"].(int))
	}
	return "-"
}

func getBytes(r *rest.Request) string {
	if r.Env["BYTES_WRITTEN"] != nil {
		return strconv.FormatInt(r.Env["BYTES_WRITTEN"].(int64), 10)
	}
	return "-"
}

func getUserAgent(r *rest.Request) string {
	if r.UserAgent() != "" {
		return r.UserAgent()
	}
	return "-"
}
