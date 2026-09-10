package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine

	requestState requestState
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		requestState: requestStateInitialized,
	}
	var toParse []byte
	bytes := make([]byte, 8)
	for request.requestState != requestStateDone {
		readBytes, err := reader.Read(bytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.requestState = requestStateDone
				break
			}
			return nil, fmt.Errorf("failed to read bytes: %w", err)
		}

		toParse = append(toParse, bytes[:readBytes]...)
		parsedBytes, err := request.parse(toParse)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bytes: %w", err)
		}

		if parsedBytes != 0 {
			toParse = toParse[parsedBytes:]
		}
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.requestState {
	case requestStateInitialized:
		parsedBytes, requestLine, err := requestLine(data)
		if err != nil {
			return 0, err
		}

		if requestLine != nil {
			r.requestState = requestStateDone
			r.RequestLine = *requestLine
			return parsedBytes, nil
		}

		return parsedBytes, nil
	case requestStateDone:
		return 0, errors.New("request already in done state")
	default:
		return 0, errors.New("unknown request state")
	}
}

func requestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, nil, nil
	}

	requestLine, err := parseRequestLineText(string(data[:idx]))
	if err != nil {
		return 0, nil, err
	}

	return idx + 2, requestLine, nil
}

func parseRequestLineText(rlt string) (*RequestLine, error) {
	requestLineParts := strings.Split(rlt, " ")
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("invalid request line text: %s", rlt)
	}

	method := requestLineParts[0]
	if len(method) == 0 {
		return nil, errors.New("invalid method length")
	}

	for _, r := range method {
		if !unicode.IsUpper(r) {
			return nil, fmt.Errorf("unknown method character: %c", r)
		}
	}

	httpVersionParts := strings.Split(requestLineParts[2], "/")
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
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersion,
	}, nil
}
