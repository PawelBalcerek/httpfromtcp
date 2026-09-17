package response

import (
	"errors"
	"fmt"
	"io"

	"github.com/PawelBalcerek/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok StatusCode = iota
	BadRequest
	InternalServerError
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case Ok:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case BadRequest:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case InternalServerError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		err = errors.New("unknown status code")
	}
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.SetContentLength(contentLen)
	h.SetConnection("close")
	h.SetContentType("text/plain")
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	for k, v := range h {
		_, err := fmt.Fprintf(w, "%s: %s\r\n", k, v)
		if err != nil {
			return fmt.Errorf("failed to write field line: %w", err)
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("failed to write end of field lines: %w", err)
	}
	return nil
}
