package response

import (
	"errors"
	"fmt"
	"io"

	"github.com/PawelBalcerek/httpfromtcp/internal/headers"
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

type writerState int

const (
	initialWriterState writerState = iota
	headersWriterState
	bodyWriterState
	doneWriterState
)

type StatusCode int

const (
	Ok StatusCode = iota
	BadRequest
	InternalServerError
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != initialWriterState {
		return errors.New("invalid writer state: not initial")
	}
	var err error
	switch statusCode {
	case Ok:
		_, err = w.writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case BadRequest:
		_, err = w.writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case InternalServerError:
		_, err = w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		err = errors.New("unknown status code")
	}
	w.writerState = headersWriterState
	return err
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != headersWriterState {
		return errors.New("invalid writer state: not headers")
	}
	for k, v := range h {
		_, err := fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
		if err != nil {
			return fmt.Errorf("failed to write field line: %w", err)
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("failed to write end of field lines: %w", err)
	}
	w.writerState = bodyWriterState
	return nil
}

func (w *Writer) WriteBody(p string) (int, error) {
	n, err := w.writeBody([]byte(p))
	if err != nil {
		return 0, err
	}
	w.writerState = doneWriterState
	return n, nil
}

func (w *Writer) writeBody(p []byte) (int, error) {
	if w.writerState != bodyWriterState {
		return 0, errors.New("invalid writer state: not body")
	}
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("failed to write payload: %w", err)
	}
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSizeN, err := w.writeBody(fmt.Appendf(nil, "%x\r\n", len(p)))
	if err != nil {
		return 0, fmt.Errorf("failed to write chunk size: %w", err)
	}
	payloadN, err := w.writeBody(p)
	if err != nil {
		return 0, fmt.Errorf("failed to write chunk: %w", err)
	}
	payloadCloseN, err := w.writeBody([]byte("\r\n"))
	if err != nil {
		return 0, fmt.Errorf("failed to write chunk delimiter: %w", err)
	}
	return chunkSizeN + payloadN + payloadCloseN, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.writeBody([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, fmt.Errorf("failed to close chunked body: %w", err)
	}
	w.writerState = doneWriterState
	return n, nil
}
