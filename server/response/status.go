package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	StatusNotFound            = "not_found"
	StatusInvalidParams       = "invalid_request_params"
	StatusUnauthorized        = "unauthorized"
	StatusConflict            = "conflict"
	StatusForbidden           = "forbidden"
	StatusTooManyRequests     = "too_many_requests"
	StatusUnprocessableEntity = "unprocessable_entity"
	StatusInternalError       = "internal_server_error"
)

var (
	OK = &Basic{Status: "ok"}

	ErrInternalServer = &Error{
		Message: "Internal error occurred. Try again later or contact service team.",
		Status:  "internal_error",
		Code:    http.StatusInternalServerError,
	}

	ErrNotFound = &Error{
		Message: "not found",
		Status:  StatusNotFound,
		Code:    http.StatusNotFound,
	}

	ErrInvalidParams = &Error{
		Message: "invalid params",
		Status:  StatusInvalidParams,
		Code:    http.StatusBadRequest,
	}

	ErrForbidden = &Error{
		Message: "forbidden",
		Status:  StatusForbidden,
		Code:    http.StatusForbidden,
	}

	ErrUnauthorized = &Error{
		Message: "Unauthorized",
		Status:  StatusUnauthorized,
		Code:    http.StatusUnauthorized,
	}

	errInvalidJSON = &Error{
		Message: "invalid json param",
		Status:  StatusInvalidParams,
		Code:    http.StatusBadRequest,
	}

	errUnknown = &Error{
		Status:  StatusInvalidParams,
		Message: "Unknown error",
		Code:    http.StatusBadRequest,
	}

	ErrTooManyRequests = &Error{
		Message: "Too many requests",
		Status:  StatusTooManyRequests,
		Code:    http.StatusTooManyRequests,
	}

	errUnprocessableEntity = &Error{
		Message: "Unprocessable",
		Status:  StatusUnprocessableEntity,
		Code:    http.StatusUnprocessableEntity,
	}

	errConflict = &Error{
		Message: "Conflict",
		Status:  StatusConflict,
		Code:    http.StatusConflict,
	}
)

type Basic struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, OK)
}

func SuccessData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func SuccessCreated(c *gin.Context) {
	c.JSON(http.StatusCreated, OK)
}

func SuccessCreatedData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func SuccessNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, OK)
}

func SuccessNotModified(c *gin.Context, _ error) {
	c.JSON(http.StatusNotModified, OK)
}

func Abort(c *gin.Context, err error) {
	var e *Error
	if errors.As(err, &e) {
		abort(c, e)
		return
	}

	AbortInternalError(c, err)
}

func AbortBadRequestData(c *gin.Context, err *Error, data interface{}) {
	err.fillMissedFields()

	res, mergeErr := mergeStruct(Basic{Message: err.Message, Status: err.Status}, data)
	if mergeErr != nil {
		AbortInternalError(c, mergeErr)
		return
	}

	c.JSON(http.StatusBadRequest, res)
}

func abort(c *gin.Context, err *Error) {
	err.fillMissedFields()

	c.Error(err.err) //nolint: errcheck

	c.AbortWithStatusJSON(err.Code, &Basic{
		Status:  err.Status,
		Message: err.Message,
		Data:    err.data,
	})
}

func AbortUnknownError(c *gin.Context, err error) {
	Abort(c, errUnknown.WithError(err))
}

func AbortTooManyRequests(c *gin.Context, err error) {
	Abort(c, ErrTooManyRequests.WithError(err))
}

func AbortUnprocessableEntity(c *gin.Context, err error) {
	Abort(c, errUnprocessableEntity.WithError(err))
}

func AbortInvalidJSON(c *gin.Context, err error) {
	Abort(c, errInvalidJSON.WithError(err))
}

func AbortInvalidParams(c *gin.Context, err error) {
	Abort(c, ErrInvalidParams.WithError(err))
}

func AbortUnauthorized(c *gin.Context, err error) {
	Abort(c, ErrUnauthorized.WithError(err))
}

func AbortForbidden(c *gin.Context, err error) {
	Abort(c, ErrForbidden.WithError(err))
}

func AbortInternalError(c *gin.Context, err error) {
	Abort(c, ErrInternalServer.WithError(err))
}

func AbortConflict(c *gin.Context, err error) {
	Abort(c, errConflict.WithError(err))
}

func mergeStruct(first, second interface{}) (map[string]interface{}, error) {
	fn := func(data interface{}) (map[string]interface{}, error) {
		var res map[string]interface{}

		str, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(str, &res)
		return res, err
	}

	item, err := fn(first)
	if err != nil {
		return nil, err
	}

	data, err := fn(second)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		item[k] = v
	}
	return item, nil
}
