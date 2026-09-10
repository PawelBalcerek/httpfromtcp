package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes: %w", err)
	}
	requestLines := strings.Split(string(bytes), "\r\n")
	if len(requestLines) < 1 {
		return nil, errors.New("request is empty")
	}
	startLine := requestLines[0]
	requestLine, err := parseRequestLine(startLine)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(startLine string) (*RequestLine, error) {
	startLineParts := strings.Split(startLine, " ")
	if len(startLineParts) != 3 {
		return nil, errors.New("invalid request start line")
	}
	method := startLineParts[0]
	if len(method) == 0 {
		return nil, errors.New("invalid method length")
	}
	for _, r := range method {
		if !unicode.IsUpper(r) {
			return nil, fmt.Errorf("unknown method character: %c", r)
		}
	}
	httpVersionParts := strings.Split(startLineParts[2], "/")
	if len(httpVersionParts) != 2 {
		return nil, errors.New("invalid http version literal")
	}
	if httpVersionParts[0] != "HTTP" {
		return nil, errors.New("invalid http version")
	}
	httpVersion := httpVersionParts[1]
	if httpVersion != "1.1" {
		return nil, errors.New("invalid http version")
	}
	return &RequestLine{
		Method:        method,
		RequestTarget: startLineParts[1],
		HttpVersion:   httpVersion,
	}, nil
}
