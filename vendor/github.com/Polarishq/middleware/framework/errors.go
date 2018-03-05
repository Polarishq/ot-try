package framework

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
)

// Error represents a error interface all swagger framework errors implement
type Error interface {
	error
	Code() int32
}

type apiError struct {
	httpCode int32
	code     int32
	message  string
}

func (a *apiError) Error() string {
	return a.message
}

func (a *apiError) Code() int32 {
	return a.code
}

func (a *apiError) HTTPCode() int32 {
	return a.httpCode
}

// NewError creates a new API error with a code and a message
func NewError(httpCode, code int32, message string, args ...interface{}) Error {
	if len(args) > 0 {
		return &apiError{httpCode, code, fmt.Sprintf(message, args...)}
	}
	return &apiError{httpCode, code, message}
}

func errorAsJSON(err Error) []byte {
	b, _ := json.Marshal(struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}{err.Code(), err.Error()})
	return b
}

// ServeError the error handler interface implemenation
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.Header().Set("Content-Type", "application/json")
	if err == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(errorAsJSON(NewError(http.StatusInternalServerError, http.StatusInternalServerError, "Unknown error")))
		return
	}
	switch e := err.(type) {
	case *errors.MethodNotAllowedError:
		rw.Header().Add("Allow", strings.Join(err.(*errors.MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(asHTTPCode(int(e.Code())))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	case Error:
		rw.WriteHeader(asHTTPCode(int(getHTTPCode(e))))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	case errors.Error:
		rw.WriteHeader(asHTTPCode(int(getHTTPCode(e))))
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(e))
		}
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		if r == nil || r.Method != "HEAD" {
			rw.Write(errorAsJSON(NewError(http.StatusInternalServerError, http.StatusInternalServerError, err.Error())))
		}
	}
}

func getHTTPCode(err Error) int32 {
	errCast, ok := err.(*apiError)
	if ok {
		return errCast.HTTPCode()
	}
	return err.Code()
}

func asHTTPCode(input int) int {
	if input >= 600 {
		return 422
	}
	return input
}
