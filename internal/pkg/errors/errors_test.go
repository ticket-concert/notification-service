package errors_test

import (
	"net/http"
	"notification-service/internal/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	// Call the function under test
	err := errors.BadRequest("Bad request")
	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, "Bad request", err.Error())
	assert.Equal(t, "Bad request", errString.Message())
	assert.Equal(t, http.StatusBadRequest, errString.Code())
}

func TestNotFound(t *testing.T) {
	// Call the function under test
	err := errors.NotFound("Not found")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, "Not found", err.Error())
	assert.Equal(t, "Not found", errString.Message())
	assert.Equal(t, http.StatusNotFound, errString.Code())
}

func TestUnprocessableEntity(t *testing.T) {
	// Call the function under test
	err := errors.UnprocessableEntity("Unprocessable entity")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, "Unprocessable entity", err.Error())
	assert.Equal(t, "Unprocessable entity", errString.Message())
	assert.Equal(t, http.StatusUnprocessableEntity, errString.Code())
}

func TestCustomError(t *testing.T) {
	// Call the function under test
	err := errors.CustomError("Custom error message", 4001, http.StatusBadRequest)

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, 4001, errString.Code())
	assert.Equal(t, "Custom error message", err.Error())
	assert.Equal(t, "Custom error message", errString.Message())
	assert.Equal(t, http.StatusBadRequest, errString.HttpCode())
}

func TestConflictError(t *testing.T) {
	// Call the function under test
	err := errors.Conflict("Conflict error message")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusConflict, errString.Code())
	assert.Equal(t, "Conflict error message", err.Error())
	assert.Equal(t, "Conflict error message", errString.Message())
}

func TestInternalServerError(t *testing.T) {
	// Call the function under test
	err := errors.InternalServerError("Internal server error message")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, errString.Code())
	assert.Equal(t, "Internal server error message", err.Error())
	assert.Equal(t, "Internal server error message", errString.Message())
}

func TestUnauthorizeError(t *testing.T) {
	// Call the function under test
	err := errors.UnauthorizedError("Unauthorize error message")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusUnauthorized, errString.Code())
	assert.Equal(t, "Unauthorize error message", err.Error())
	assert.Equal(t, "Unauthorize error message", errString.Message())
}

func TestForbiddenError(t *testing.T) {
	// Call the function under test
	err := errors.ForbiddenError("Forbidden error message")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusForbidden, errString.Code())
	assert.Equal(t, "Forbidden error message", err.Error())
	assert.Equal(t, "Forbidden error message", errString.Message())
}

func TestTooManyRequestError(t *testing.T) {
	// Call the function under test
	err := errors.TooManyRequest("Too many request error message")

	errString, _ := err.(*errors.ErrorString)
	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusTooManyRequests, errString.Code())
	assert.Equal(t, "Too many request error message", err.Error())
	assert.Equal(t, "Too many request error message", errString.Message())
}
