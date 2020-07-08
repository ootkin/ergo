package ergo

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Application error codes
const (
	ECONFLICT = "conflict"  // Action cannot be performed
	EINTERNAL = "internal"  // Internal error
	EINVALID  = "invalid"   // Validation failed
	ENOTFOUND = "not_found" // Entity does not exists
)

// Error defines a standard application error
// Code is a Machine-readable error code
// Message is a Human-readable message
// Op is the logical operation that has generated the error
// Err is the error generated
type Error struct {
	Code    string
	Message string
	Op      string
	Err     error
}

// JSON Error defines the error to send to client
type JSONError struct {
	Code       string `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// Error returns the string representation of the error message.
func (err *Error) Error() string {
	var buffer bytes.Buffer

	// Print the current operation in our stack, if any
	if err.Op != "" {
		_, _ = fmt.Fprintf(&buffer, "%s: ", err.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if err.Err != nil {
		buffer.WriteString(err.Err.Error())
	} else {
		if err.Code != "" {
			_, _ = fmt.Fprintf(&buffer, "<%s>", err.Code)
		}
		buffer.WriteString(err.Message)
	}
	return buffer.String()
}

// ErrorCode returns the code of the root error, if available.
// Otherwise returns EINTERNAL.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	} else if e, isCustomError := err.(*Error); isCustomError && e.Code != "" {
		return e.Code
	} else if isCustomError && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return EINTERNAL
}

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise returns a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	} else if e, isCustomError := err.(*Error); isCustomError && e.Message != "" {
		return e.Message
	} else if isCustomError && e.Err != nil {
		return ErrorMessage(e.Err)
	} else if isCustomError && e.Code != "" {
		// If the message is not present, try to infer it from the Code
		switch e.Code {
		case ECONFLICT:
			return "Conflict error."
		case EINTERNAL:
			return "An internal error has occurred."
		case EINVALID:
			return "Bad request."
		case ENOTFOUND:
			return "Resource not found."
		}
	}
	return "An internal error has occurred."
}

// ErrorStatusCode returns the status code of the http request.
// Otherwise returns a 500 (internal server error)
func ErrorStatusCode(err error) int {
	if e, isCustomError := err.(*Error); isCustomError && e.Code != "" {
		switch e.Code {
		case ECONFLICT:
			return http.StatusConflict
		case EINTERNAL:
			return http.StatusInternalServerError
		case EINVALID:
			return http.StatusBadRequest
		case ENOTFOUND:
			return http.StatusNotFound
		}
	} else if isCustomError && e.Err != nil {
		return ErrorStatusCode(e.Err)
	}
	// Fallback
	return http.StatusInternalServerError
}

// Format error will return a Json to be sent to the client describing the error
func FormatError(err error) JSONError {
	return JSONError{
		Code:       ErrorCode(err),
		StatusCode: ErrorStatusCode(err),
		Message:    ErrorMessage(err),
	}
}

// HandleError will return a Json representation of the error and log the error
func HandleError(err error) (int, JSONError) {
	log.Println(err.Error())
	return ErrorStatusCode(err), FormatError(err)
}
