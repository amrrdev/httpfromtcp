package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

var validHeaderName = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+$`)
var rn = []byte("\r\n")

func ParseHeader(fields []byte) (string, string, error) {
	parts := bytes.SplitN(fields, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])
	if len(bytes.TrimSpace(name)) == 0 {
		return "", "", fmt.Errorf("empty field name")
	}

	for _, b := range name {
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			return "", "", fmt.Errorf("malformed field name")
		}
	}

	if !validHeaderName.MatchString(string(name)) {
		return "", "", fmt.Errorf("malformed field name (invalid characters)")
	}

	return strings.ToLower(string(name)), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0

	for {
		idx := bytes.Index(data, rn)
		if idx == -1 {
			break
		}

		// EMPTY HEADER
		if idx == 0 {
			read += len(rn)
			break
		}

		name, value, err := ParseHeader(data[:idx])
		if err != nil {
			return 0, false, err
		}

		h[name] = value
		read += idx + len(rn)
		data = data[idx+len(rn):]

	}

	return read, true, nil
}
