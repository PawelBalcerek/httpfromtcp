package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) ParseSingleHeader(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	header := string(data[:idx])
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) != 2 {
		return 0, false, fmt.Errorf("invalid header line: %s", header)
	}

	name := headerParts[0]
	if strings.ContainsAny(name, " \t") {
		return 0, false, fmt.Errorf("invalid header name: %s (contains whitespace characters)", name)
	}

	value := strings.TrimSpace(headerParts[1])
	h.SetHeader(name, value)

	return idx + 2, false, nil
}

func (h Headers) SetHeader(name, value string) {
	h[name] = value
}
