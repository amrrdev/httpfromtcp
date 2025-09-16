package headers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaderParser(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nContent-Length: 100\r\nFooFoo: BarBar\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers["Host"])
	require.Equal(t, "BarBar", headers["FooFoo"])
	require.Equal(t, "", headers["MissingKey"])
	require.True(t, done)
	require.Equal(t, 60, n)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
}

func TestHeaderParserEmptyHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, 2, n) // Just \r\n
	require.True(t, done)
	require.Len(t, headers, 0)
}

func TestHeaderParserSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-Type: application/json\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "application/json", headers["Content-Type"])
	// require.Equal(t, 32, n)
	require.True(t, done)
}

func TestHeaderParserValueTrimming(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Authorization:   Bearer token123   \r\nHost:  example.com  \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "Bearer token123", headers["Authorization"])
	require.Equal(t, "example.com", headers["Host"])
	require.True(t, done)
}

func TestHeaderParserEmptyValue(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Custom-Header:\r\nAnother-Header:   \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "", headers["X-Custom-Header"])
	require.Equal(t, "", headers["Another-Header"])
	require.True(t, done)
}

func TestHeaderParserMalformedHeaders(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "no colon",
			data: []byte("InvalidHeader\r\n\r\n"),
		},
		{
			name: "multiple colons but valid",
			data: []byte("Time: 12:34:56\r\n\r\n"), // This should be valid
		},
		{
			name: "header name with trailing space",
			data: []byte("Host : localhost\r\n\r\n"),
		},
		{
			name: "header name with tab",
			data: []byte("Host\t: localhost\r\n\r\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := NewHeaders()
			n, done, err := headers.Parse(tc.data)
			if tc.name == "multiple colons but valid" {
				require.NoError(t, err)
				require.Equal(t, "12:34:56", headers["Time"])
				require.True(t, done)
			} else {
				require.Error(t, err)
				require.Equal(t, 0, n)
				require.False(t, done)
			}
		})
	}
}

func TestHeaderParserIncompleteData(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost\r\nContent-Length: 100")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "localhost", headers["Host"])
	require.True(t, done)   // No error, but incomplete
	require.Equal(t, 17, n) // Only "Host: localhost\r\n"
}

func TestHeaderParserNoTermination(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, 0, n) // No \r\n found, nothing parsed
	require.True(t, done)
	require.Len(t, headers, 0)
}

func TestHeaderParserCasePreservation(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-TYPE: text/html\r\nX-CUSTOM-header: value\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "text/html", headers["Content-TYPE"])
	require.Equal(t, "value", headers["X-CUSTOM-header"])
	require.True(t, done)
}

func TestHeaderParserSpecialCharacters(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Special: !@#$%^&*()_+-=[]{}|;:,.<>?\r\nX-Unicode: héllo wørld\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "!@#$%^&*()_+-=[]{}|;:,.<>?", headers["X-Special"])
	require.Equal(t, "héllo wørld", headers["X-Unicode"])
	require.True(t, done)
}

func TestHeaderParserLongHeaders(t *testing.T) {
	longValue := make([]byte, 1000)
	for i := range longValue {
		longValue[i] = 'a'
	}

	headers := NewHeaders()
	data := []byte("X-Long-Header: " + string(longValue) + "\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, string(longValue), headers["X-Long-Header"])
	require.True(t, done)
}

func TestHeaderParserDuplicateHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Set-Cookie: first=value1\r\nSet-Cookie: second=value2\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	// Last one wins (typical map behavior)
	require.Equal(t, "second=value2", headers["Set-Cookie"])
	require.True(t, done)
}

func TestHeaderParserOnlyLineFeed(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost\n\n") // \n instead of \r\n
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, 0, n) // No \r\n found
	require.True(t, done)
	require.Len(t, headers, 0)
}

func TestHeaderParserMixedLineEndings(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost\r\nContent-Type: text/html\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "localhost", headers["Host"])
	// Should stop after first header since second line doesn't end with \r\n
	require.Equal(t, 43, n) // "Host: localhost\r\n"
	require.True(t, done)
}

func TestHeaderParserEmptyFieldName(t *testing.T) {
	headers := NewHeaders()
	data := []byte(": value\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
}

func TestHeaderParserTabsInValue(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Tabs:\t\tvalue\twith\ttabs\t\t\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "value\twith\ttabs", headers["X-Tabs"]) // TrimSpace should remove leading/trailing tabs
	require.True(t, done)
}
