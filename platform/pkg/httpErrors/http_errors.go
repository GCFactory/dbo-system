package httpErrors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	ErrBadRequest       = "Bad request"
	ErrWrongCredentials = "Wrong Credentials"
	ErrNotFound         = "Not Found"
	ErrUnauthorized     = "Unauthorized"
	ErrForbidden        = "Forbidden"
	ErrBadQueryParams   = "Invalid query params"
)

var (
	BadRequest          = errors.New("Bad request")
	WrongCredentials    = errors.New("Wrong Credentials")
	NotFound            = errors.New("Not Found")
	Unauthorized        = errors.New("Unauthorized")
	Forbidden           = errors.New("Forbidden")
	PermissionDenied    = errors.New("Permission Denied")
	ExpiredCSRFError    = errors.New("Expired CSRF token")
	WrongCSRFToken      = errors.New("Wrong CSRF token")
	CSRFNotPresented    = errors.New("CSRF not presented")
	NotRequiredFields   = errors.New("No such required fields")
	BadQueryParams      = errors.New("Invalid query params")
	InternalServerError = errors.New("Internal Server Error")
	RequestTimeoutError = errors.New("Request Timeout")
	InvalidJWTToken     = errors.New("Invalid JWT token")
	InvalidJWTClaims    = errors.New("Invalid JWT claims")
	NoCookie            = errors.New("not found cookie header")
)

// Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

// Rest error struct
type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

// Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// New Rest Error
func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New Rest Error With Message
func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// New Internal Server Error
func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

// Parser of error string messages returns RestError
func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, NotFound.Error(), err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, RequestTimeoutError.Error(), err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	case strings.Contains(err.Error(), "Field validation"):
		return parseValidatorError(err)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
	case strings.Contains(err.Error(), "UUID"):
		return NewRestError(http.StatusBadRequest, err.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "cookie"):
		return NewRestError(http.StatusUnauthorized, Unauthorized.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewRestError(http.StatusUnauthorized, Unauthorized.Error(), err)
	case strings.Contains(strings.ToLower(err.Error()), "bcrypt"):
		return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
	default:
		if restErr, ok := err.(RestErr); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func parseSqlErrors(err error) RestErr {
	//if strings.Contains(err.Error(), "23505") {
	//	return NewRestError(http.StatusBadRequest, ExistsEmailError.Error(), err)
	//}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

func parseValidatorError(err error) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, "Invalid password, min length 6", err)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, "Invalid email", err)
	}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

// Error response
func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
