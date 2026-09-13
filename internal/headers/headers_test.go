package headers_test

import (
	"testing"

	"github.com/PawelBalcerek/httpfromtcp/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	h := headers.NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := h.ParseSingleHeader(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespaces
	h = headers.NewHeaders()
	data = []byte("Host:   localhost:42069  \r\n\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, 27, n)
	assert.False(t, done)

	// Test: Valid two headers with existing header
	data = []byte("Authorization: Bearer token\r\nHost: localhost:42069\r\n\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "localhost:42069", h["host"])
	assert.Equal(t, "Bearer token", h["authorization"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// Test: Valid done
	h = headers.NewHeaders()
	data = []byte("\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Duplicated header
	h = headers.NewHeaders()
	data = []byte("X-Custom-Header: value1\r\n\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "value1", h["x-custom-header"])
	assert.Equal(t, 25, n)
	assert.False(t, done)
	data = []byte("X-Custom-Header: value2\r\n\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Equal(t, "value1, value2", h["x-custom-header"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	h = headers.NewHeaders()
	data = []byte(" Host : localhost:42069\r\n\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Non-compliant header field-name
	h = headers.NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = h.ParseSingleHeader(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
