package custom_error

type AppError struct {
	Code    int
	Message string
}

// Implement the error interface
func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
