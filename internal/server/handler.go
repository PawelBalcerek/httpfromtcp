package server

import (
	"fmt"
	"io"

	"github.com/PawelBalcerek/httpfromtcp/internal/request"
	"github.com/PawelBalcerek/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he HandlerError) WriteError(w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return fmt.Errorf("failed to write status line: %w", err)
	}

	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(he.Message)))
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	_, err = w.Write([]byte(he.Message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
