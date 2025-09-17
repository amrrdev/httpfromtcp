package headers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaderParser(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nContent-Length: 100\r\nFooFoo: BarBar\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, "100", headers["content-length"])
	require.Equal(t, "BarBar", headers["foofoo"])
	require.Equal(t, "", headers["MissingKey"])
	require.True(t, done)
	require.Equal(t, 62, n)
}

func TestHeaderParserDuplicateHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Set-Cookie: first=value1\r\nSet-Cookie: second=value2\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	// Multiple headers should be combined with commas per RFC
	require.Equal(t, "first=value1, second=value2", headers["set-cookie"])
	require.True(t, done)
}

func TestHeaderParserIncompleteData(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost\r\nContent-Length: 100")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "localhost", headers["host"])
	require.False(t, done)
	require.Equal(t, 17, n)
}

func TestHeaderParserNoTermination(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	require.Len(t, headers, 0)
}

func TestHeaderParserMultipleDuplicateHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Accept: text/html\r\nAccept: application/json\r\nAccept: */*\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "text/html, application/json, */*", headers["accept"])
	require.True(t, done)
}

func TestHeaderParserDuplicateHeadersWithSpacing(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Custom:   value1   \r\nX-Custom:   value2   \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "value1, value2", headers["x-custom"])
	require.True(t, done)
}

func TestHeaderParserDuplicateHeadersDifferentCases(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-Type: text/html\r\ncontent-type: application/json\r\nCONTENT-TYPE: text/plain\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "text/html, application/json, text/plain", headers["content-type"])
	require.True(t, done)
}

func TestHeaderParserEmptyValueInDuplicates(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Test: value1\r\nX-Test:\r\nX-Test: value2\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "value1, , value2", headers["x-test"])
	require.True(t, done)
}

func TestHeaderParserSingleHeaderNoDuplicate(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Authorization: Bearer token123\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "Bearer token123", headers["authorization"])
	require.True(t, done)
}

func TestHeaderParserManyDuplicates(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Via: proxy1\r\nVia: proxy2\r\nVia: proxy3\r\nVia: proxy4\r\nVia: proxy5\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "proxy1, proxy2, proxy3, proxy4, proxy5", headers["via"])
	require.True(t, done)
}

func TestHeaderParserDuplicatesWithOtherHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: example.com\r\nAccept: text/html\r\nAccept: application/json\r\nContent-Length: 100\r\nAccept: */*\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "example.com", headers["host"])
	require.Equal(t, "text/html, application/json, */*", headers["accept"])
	require.Equal(t, "100", headers["content-length"])
	require.True(t, done)
	require.Len(t, headers, 3)
}

func TestHeaderParserDuplicateWithCommaInValue(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Accept-Language: en-US,en;q=0.9\r\nAccept-Language: fr;q=0.8\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "en-US,en;q=0.9, fr;q=0.8", headers["accept-language"])
	require.True(t, done)
}

func TestHeaderParserDuplicateEmptyAndNonEmpty(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Empty:\r\nX-Empty: not-empty\r\nX-Empty:\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, ", not-empty, ", headers["x-empty"])
	require.True(t, done)
}
