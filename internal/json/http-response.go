package json

import (
	"fmt"
	"net/http"

	"github.com/ricohartono89/base-api/errs"
)

// HTTPResponse ...
type HTTPResponse struct {
	status  int
	data    interface{}
	err     error
	errCode int
	message string
}

// SetOk ...
func (res HTTPResponse) SetOk(data interface{}) HTTPResponse {
	res.status = http.StatusOK
	res.data = data
	return res
}

// SetOkWithStatus ...
func (res HTTPResponse) SetOkWithStatus(status int, data interface{}) HTTPResponse {
	res.status = status
	res.data = data
	return res
}

// SetError ...
func (res HTTPResponse) SetError(err error, errCode int, message string) HTTPResponse {
	res.status = http.StatusInternalServerError
	res.err = err
	res.errCode = errCode
	res.message = message
	return res
}

// SetErrorWithStatus ...
func (res HTTPResponse) SetErrorWithStatus(status int, err error, errCode int, message string) HTTPResponse {
	res.status = status
	res.err = err
	res.errCode = errCode
	res.message = message
	return res
}

// ImportJSONWebError ...
func (res HTTPResponse) ImportJSONWebError(err *errs.JSONWebError) HTTPResponse {
	res.status = err.Status
	res.err = err.Error
	res.errCode = err.Code
	res.message = err.Message
	return res
}

// HasError ...
func (res HTTPResponse) HasError() bool {
	return res.err != nil
}

// GetData ...
func (res HTTPResponse) GetData() interface{} {
	return res.data
}

// GetError ...
func (res HTTPResponse) GetError() error {
	return res.err
}

// GetStatus ...
func (res HTTPResponse) GetStatus() int {
	if res.status != 0 {
		return res.status
	}
	return http.StatusInternalServerError
}

// GetErrCode ...
func (res HTTPResponse) GetErrCode() int {
	if res.errCode != 0 {
		return res.errCode
	}
	return errs.Unknown
}

// GetErrorMessage get error message from message or error object
func (res HTTPResponse) GetErrorMessage() string {
	if res.message != "" {
		return res.message
	}
	return res.err.Error()
}

// GetErrorMessageVerbose get full string with error code, message and error object
func (res HTTPResponse) GetErrorMessageVerbose() string {
	return fmt.Sprintf("Error Code: %d, Message: %s. Detail: %s", res.errCode, res.message, res.err.Error())
}
