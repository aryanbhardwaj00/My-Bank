package customerrors

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("Resource Not Found")
	ErrInvalidInput    = errors.New("Invalid Input Provided")
	ErrConnection      = errors.New("Connection Failed")
	ErrInternalServer  = errors.New("Internal Server Error")
	ErrResourceDeleted = errors.New("Requested resource Deleted")
	ErrNoRows          = errors.New("No Rows Affected")
)
