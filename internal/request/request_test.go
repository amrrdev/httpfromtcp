package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestLineParser(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.HttpMethod)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.HttpMethod)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

	_, err = RequestFromReader(strings.NewReader("GET /amr HTTP/1.2\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)
}
