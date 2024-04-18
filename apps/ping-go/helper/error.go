package helper

import "fmt"

func NewError(message string, code int, errors interface{}) error {
	return &ResponseBody{Code: code, Message: message, Errors: errors}
}

// Error implements error.
func (r *ResponseBody) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Errors: %v \n", r.Code, r.Message, r.Errors)
}
