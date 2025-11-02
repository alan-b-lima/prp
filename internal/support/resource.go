package support

import (
	"encoding/json"
	"net/http"

	"github.com/alan-b-lima/prp/pkg/errors"
)

type RouteFunc func(w http.ResponseWriter, r *http.Request) error

func (fn RouteFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}

	if err, ok := errors.AsType[*errors.Error](err); ok {
		writeJsonError(w, err, toStatusCode(err.Kind))
		return
	}

	writeJsonError(w, err, http.StatusInternalServerError)
}

func writeJsonError(w http.ResponseWriter, err error, status int) {
	body, e := json.Marshal(err)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body)
}

var StatusCodes = map[errors.Kind]int{
	errors.InvalidInput:       http.StatusBadRequest,
	errors.Unauthorized:       http.StatusUnauthorized,
	errors.PreconditionFailed: http.StatusPreconditionFailed,
	errors.NotFound:           http.StatusNotFound,
	errors.Conflict:           http.StatusConflict,

	errors.Internal:    http.StatusInternalServerError,
	errors.Unavailable: http.StatusServiceUnavailable,
	errors.Timeout:     http.StatusRequestTimeout,
	errors.BadGateway:  http.StatusBadGateway,
}

func toStatusCode(kind errors.Kind) int {
	if status, in := StatusCodes[kind]; in {
		return status
	}

	return http.StatusInternalServerError
}
