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
	require.Equal(t, "localhost:42069", headers["host"]) // Should be lowercase now
	require.Equal(t, "BarBar", headers["foofoo"])        // Should be lowercase now
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
	require.Equal(t, "application/json", headers["content-type"]) // Should be lowercase now
	require.True(t, done)
}

func TestHeaderParserValueTrimming(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Authorization:   Bearer token123   \r\nHost:  example.com  \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "Bearer token123", headers["authorization"]) // Should be lowercase now
	require.Equal(t, "example.com", headers["host"])              // Should be lowercase now
	require.True(t, done)
}

func TestHeaderParserEmptyValue(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Custom-Header:\r\nAnother-Header:   \r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "", headers["x-custom-header"]) // Should be lowercase now
	require.Equal(t, "", headers["another-header"])  // Should be lowercase now
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
				require.Equal(t, "12:34:56", headers["time"]) // Should be lowercase now
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
	require.Equal(t, "localhost", headers["host"]) // Should be lowercase now
	require.True(t, done)                          // No error, but incomplete
	require.Equal(t, 17, n)                        // Only "Host: localhost\r\n"
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

// UPDATED: This test now expects lowercase keys
func TestHeaderParserCaseInsensitivity(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-TYPE: text/html\r\nX-CUSTOM-header: value\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "text/html", headers["content-type"]) // Should be lowercase
	require.Equal(t, "value", headers["x-custom-header"])  // Should be lowercase
	require.True(t, done)

	// Verify original casing doesn't exist
	require.Equal(t, "", headers["Content-TYPE"])
	require.Equal(t, "", headers["X-CUSTOM-header"])
}

func TestHeaderParserSpecialCharacters(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-Special: !@#$%^&*()_+-=[]{}|;:,.<>?\r\nX-Unicode: héllo wørld\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Equal(t, "!@#$%^&*()_+-=[]{}|;:,.<>?", headers["x-special"]) // Should be lowercase now
	require.Equal(t, "héllo wørld", headers["x-unicode"])                // Should be lowercase now
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
	require.Equal(t, string(longValue), headers["x-long-header"]) // Should be lowercase now
	require.True(t, done)
}

func TestHeaderParserDuplicateHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Set-Cookie: first=value1\r\nSet-Cookie: second=value2\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	// Last one wins (typical map behavior)
	require.Equal(t, "second=value2", headers["set-cookie"]) // Should be lowercase now
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
	require.Equal(t, "localhost", headers["host"]) // Should be lowercase now
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
	require.Equal(t, "value\twith\ttabs", headers["x-tabs"]) // Should be lowercase now, TrimSpace should remove leading/trailing tabs
	require.True(t, done)
}

// NEW TESTS FOR THE REQUIREMENTS

func TestHeaderParserCaseInsensitivityRequirement(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "All uppercase",
			input:    "CONTENT-LENGTH: 100\r\n\r\n",
			expected: "content-length",
		},
		{
			name:     "Mixed case",
			input:    "Content-Length: 100\r\n\r\n",
			expected: "content-length",
		},
		{
			name:     "All lowercase",
			input:    "content-length: 100\r\n\r\n",
			expected: "content-length",
		},
		{
			name:     "Random case",
			input:    "CoNtEnT-LeNgTh: 100\r\n\r\n",
			expected: "content-length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := NewHeaders()
			_, done, err := headers.Parse([]byte(tc.input))
			require.NoError(t, err)
			require.True(t, done)
			require.Equal(t, "100", headers[tc.expected])
			// Verify we have exactly one key and it's lowercase
			require.Len(t, headers, 1)
		})
	}
}

func TestHeaderParserInvalidCharacters(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "copyright symbol in header name",
			data: []byte("H©st: localhost:42069\r\n\r\n"),
		},
		{
			name: "emoji in header name",
			data: []byte("X-Custom-🚀: rocket\r\n\r\n"),
		},
		{
			name: "parentheses in header name (invalid)",
			data: []byte("X-Custom(test): value\r\n\r\n"),
		},
		{
			name: "comma in header name (invalid)",
			data: []byte("X-Custom,test: value\r\n\r\n"),
		},
		{
			name: "semicolon in header name (invalid)",
			data: []byte("X-Custom;test: value\r\n\r\n"),
		},
		{
			name: "equals in header name (invalid)",
			data: []byte("X-Custom=test: value\r\n\r\n"),
		},
		{
			name: "square brackets in header name (invalid)",
			data: []byte("X-Custom[test]: value\r\n\r\n"),
		},
		{
			name: "curly braces in header name (invalid)",
			data: []byte("X-Custom{test}: value\r\n\r\n"),
		},
		{
			name: "quotes in header name (invalid)",
			data: []byte("X-Custom\"test\": value\r\n\r\n"),
		},
		{
			name: "backslash in header name (invalid)",
			data: []byte("X-Custom\\test: value\r\n\r\n"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := NewHeaders()
			n, done, err := headers.Parse(tc.data)
			require.Error(t, err)
			require.Equal(t, 0, n)
			require.False(t, done)
			require.Contains(t, err.Error(), "malformed field name")
		})
	}
}

func TestHeaderParserValidSpecialCharacters(t *testing.T) {
	// Test all valid special characters from the RFC: ! # $ % & ' * + - . ^ _ ` | ~
	testCases := []struct {
		name   string
		header string
	}{
		{"exclamation", "X-Test!: value"},
		{"hash", "X-Test#: value"},
		{"dollar", "X-Test$: value"},
		{"percent", "X-Test%: value"},
		{"ampersand", "X-Test&: value"},
		{"apostrophe", "X-Test': value"},
		{"asterisk", "X-Test*: value"},
		{"plus", "X-Test+: value"},
		{"dash", "X-Test-: value"},
		{"dot", "X-Test.: value"},
		{"caret", "X-Test^: value"},
		{"underscore", "X-Test_: value"},
		{"backtick", "X-Test`: value"},
		{"pipe", "X-Test|: value"},
		{"tilde", "X-Test~: value"},
		{"mixed valid chars", "X-Test!#$%&'*+-.^_`|~123: value"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(tc.header + "\r\n\r\n")
			_, done, err := headers.Parse(data)
			require.NoError(t, err, "Should accept valid header: %s", tc.header)
			require.True(t, done)
			require.Len(t, headers, 1)
		})
	}
}

func TestHeaderParserCaseInsensitivityWithValidChars(t *testing.T) {
	headers := NewHeaders()
	data := []byte("X-CUSTOM-Test123!#$: value\r\nContent-TYPE: text/html\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.True(t, done)

	// Both should be stored with lowercase keys
	require.Equal(t, "value", headers["x-custom-test123!#$"])
	require.Equal(t, "text/html", headers["content-type"])

	// Verify original case doesn't exist
	require.Equal(t, "", headers["X-CUSTOM-Test123!#$"])
	require.Equal(t, "", headers["Content-TYPE"])
}
