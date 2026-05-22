package rag

type validationError struct {
	message string
}

func newValidationError(message string) validationError {
	return validationError{message: message}
}

func (e validationError) Error() string {
	return e.message
}
