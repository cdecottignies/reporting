package client

import "errors"

var ErrNotFound = errors.New("not found")

type ErrorMessage struct {
	ErrorPayload `json:"error"`
}

type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

func (e *ErrorMessage) Error() string {
	return e.Message
}
