package ergo

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := &Error{
		Code:    EINVALID,
		Message: "resource not found",
		Op:      "operation.test",
	}
	actual := err.Error()
	expected := "operation.test: <invalid>resource not found"
	assert.Equal(t, expected, actual)

	err = &Error{
		Err: errors.New("error message"),
	}
	actual = err.Error()
	expected = "error message"
	assert.Equal(t, expected, actual)
}

func TestErrorCode(t *testing.T) {
	// Test with error as nil
	actual := ErrorCode(nil)
	assert.Equal(t, "", actual)

	// Test with normal error
	error := errors.New("some error")
	actual = ErrorCode(error)
	assert.Equal(t, EINTERNAL, actual)

	// Test with Code in Error
	error = &Error{
		Code:    EINVALID,
		Message: "message",
		Op:      "operation",
		Err:     nil,
	}
	actual = ErrorCode(error)
	assert.Equal(t, EINVALID, actual)

	// Test without Code in Error but with error
	error = &Error{
		Code:    "",
		Message: "message",
		Op:      "operation",
		Err:     errors.New("some error"),
	}
	actual = ErrorCode(error)
	assert.Equal(t, EINTERNAL, actual)
}

func TestErrorMessage(t *testing.T) {
	// Test with error as nil
	actual := ErrorMessage(nil)
	assert.Equal(t, "", actual)

	// Test with normal error
	error := errors.New("some error")
	actual = ErrorMessage(error)
	assert.Equal(t, "An internal error has occurred.", actual)

	// Test with Message in Error
	error = &Error{
		Code:    "",
		Message: "error message",
		Op:      "",
		Err:     nil,
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "error message", actual)

	// Test with No Message in Error
	error = &Error{
		Code:    "",
		Message: "",
		Op:      "",
		Err:     nil,
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "An internal error has occurred.", actual)

	// Test with an error in Error
	error = &Error{
		Code:    "",
		Message: "",
		Op:      "",
		Err:     errors.New("some error"),
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "An internal error has occurred.", actual)

	// Infer message from Error Code
	error = &Error{
		Code: ECONFLICT,
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "Conflict error.", actual)

	error = &Error{
		Code: EINTERNAL,
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "An internal error has occurred.", actual)

	error = &Error{
		Code: EINVALID,
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "Bad request.", actual)

	error = &Error{
		Code: ENOTFOUND,
	}
	actual = ErrorMessage(error)
	assert.Equal(t, "Resource not found.", actual)
}

func TestErrorStatusCode(t *testing.T) {
	// Test with nil error
	actual := ErrorStatusCode(nil)
	assert.Equal(t, http.StatusInternalServerError, actual)

	// Test with normal error
	error := errors.New("some error")
	actual = ErrorStatusCode(error)
	assert.Equal(t, http.StatusInternalServerError, actual)

	// Test with Code in Error
	error = &Error{
		Code: ENOTFOUND,
	}
	actual = ErrorStatusCode(error)
	assert.Equal(t, http.StatusNotFound, actual)
	error = &Error{
		Code: EINTERNAL,
	}
	actual = ErrorStatusCode(error)
	assert.Equal(t, http.StatusInternalServerError, actual)

	// Test with error in Error
	error = &Error{
		Err: &Error{
			Code: ECONFLICT,
		},
	}
	actual = ErrorStatusCode(error)
	assert.Equal(t, http.StatusConflict, actual)
}

func TestFormatError(t *testing.T) {
	// Test with nil
	actual := FormatError(nil)
	expected := JSONError{
		Code:       "",
		StatusCode: http.StatusInternalServerError,
		Message:    "",
	}
	assert.Equal(t, expected, actual)

	// Test with Error
	error := &Error{
		Code:    EINVALID,
		Message: "message",
	}
	expected = JSONError{
		Code:       EINVALID,
		StatusCode: http.StatusBadRequest,
		Message:    "message",
	}
	actual = FormatError(error)
	assert.Equal(t, expected, actual)

	// Test without Message in error, should infer it from Code
	error = &Error{
		Code: EINVALID,
	}
	expected = JSONError{
		Code:       EINVALID,
		StatusCode: http.StatusBadRequest,
		Message:    "Bad request.",
	}
	actual = FormatError(error)
	assert.Equal(t, expected, actual)
}

func TestHandleError(t *testing.T) {
	error := &Error{
		Code:    EINVALID,
		Message: "custom message",
	}
	expectedHttpStatus := http.StatusBadRequest
	expectedJsonError := JSONError{
		Code:       EINVALID,
		StatusCode: http.StatusBadRequest,
		Message:    "custom message",
	}
	actualHttpStatus, actualJsonError := HandleError(error)
	assert.Equal(t, expectedHttpStatus, actualHttpStatus)
	assert.Equal(t, expectedJsonError, actualJsonError)
}
