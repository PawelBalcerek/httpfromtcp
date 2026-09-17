package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var fieldNameRegex = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`)

const (
	ContentLength = "Content-Length"
	Connection    = "Connection"
	ContentType   = "Content-Type"
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
	if !fieldNameRegex.MatchString(name) {
		return 0, false, fmt.Errorf("non-compliant header field-name: %s", name)
	}

	value := strings.TrimSpace(headerParts[1])
	h.SetHeader(name, value)

	return idx + 2, false, nil
}

func (h Headers) SetContentLength(cl int) {
	h.SetHeader(ContentLength, strconv.Itoa(cl))
}

func (h Headers) SetConnection(v string) {
	h.SetHeader(Connection, v)
}

func (h Headers) SetContentType(v string) {
	h.SetHeader(ContentType, v)
}

func (h Headers) SetHeader(name, value string) {
	lName := strings.ToLower(name)
	if v, ok := h[lName]; ok {
		h[lName] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h[lName] = value
	}
}

func (h Headers) GetContentLength() (value int, ok bool, err error) {
	if v, ok := h.GetHeader(ContentLength); ok {
		value, err := strconv.Atoi(v)
		if err != nil {
			return 0, false, fmt.Errorf("failed to convert %s header: %w", ContentLength, err)
		}
		return value, true, nil
	} else {
		return 0, false, nil
	}
}

func (h Headers) GetHeader(name string) (value string, ok bool) {
	lName := strings.ToLower(name)
	value, ok = h[lName]
	return value, ok
}
