package exception

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Error struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
}

func NewError() *Error {
	return &Error{}
}

func (e *Error) ProductNotFoundError() *Error {
	return &Error{"0001",
		"Product not foud",
		time.Now().UTC().Format(time.RFC3339),
		http.StatusNotFound}
}

func (e *Error) ProductNotFoundErrorWithMessage(message string) *Error {
	return &Error{"0001",
		message,
		time.Now().UTC().Format(time.RFC3339),
		http.StatusNotFound}
}

func (e *Error) ToJSON(w io.Writer) error {
	j := json.NewEncoder(w)
	return j.Encode(e)
}
