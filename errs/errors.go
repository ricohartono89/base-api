package errs

import (
	"fmt"
	"strconv"
)

// Unknown ...
const Unknown = 1

// ServiceID ...
const ServiceID = "1"

// JSONWebError contains error information
type JSONWebError struct {
	Error   error
	Status  int
	Code    int
	Message string
}

// String get string with error code and message
func (e JSONWebError) String() string {
	var message string
	if e.Message == "" {
		message = e.Error.Error()
	} else {
		message = e.Message
	}
	return fmt.Sprintf("Error Code: %d, Message: %s", e.Code, message)
}

// ErrorMessage get error message from message or error object
func (e JSONWebError) ErrorMessage() string {
	if e.Message == "" {
		return e.Error.Error()
	}
	return e.Message
}

// StringWithError get full string with error code, message and error object
func (e JSONWebError) StringWithError() string {
	return fmt.Sprintf("Error Code: %d, Message: %s. Detail: %s", e.Code, e.Message, e.Error.Error())
}

// BuildJSONWebError build json web error
// Deprecated: correspond to JSONWebError, use BuildErrorCode instead
func BuildJSONWebError(prefixCode string, suffixCode string, err error, status int) *JSONWebError {
	jsonWebError := &JSONWebError{
		Code:    BuildErrorCode(prefixCode, suffixCode),
		Status:  status,
		Message: err.Error(),
		Error:   err,
	}
	return jsonWebError
}

// BuildErrorCode ...
func BuildErrorCode(prefixCode string, suffixCode string) int {
	var code int
	errorString := ServiceID + prefixCode + suffixCode
	code, e := strconv.Atoi(errorString)
	if e != nil {
		code = Unknown
	}

	return code
}
