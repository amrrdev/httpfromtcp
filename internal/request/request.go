package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpMethod    string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        string
}

var (
	SEPARATOR                       = "\r\n"
	ERROR_MALFORM_REQUEST_LINE      = fmt.Errorf("malformed request-line")
	ERROR_UNSUPPORTTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll"), err)
	}

	rl, _, err := parseRequestLine(string(data))

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, err
}

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, b, ERROR_MALFORM_REQUEST_LINE
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPARATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, b, ERROR_MALFORM_REQUEST_LINE
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMsg, ERROR_MALFORM_REQUEST_LINE
	}

	rl := &RequestLine{
		HttpMethod:    parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpParts[1],
	}

	return rl, restOfMsg, nil
}
