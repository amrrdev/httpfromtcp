package headers

import (
	"bytes"
	"errors"
	"fmt"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

/*
Make sure there’s no space between the header field name and the colon.
Wrong:  Host : amr.com
Right:   Host: amr.com

Example:
  Host: amrmubarak.com\r\n
  Content-Length: 902\r\n
  \r\n   // end of headers
*/

// Look for "\r\n\r\n" — that marks the end of the HTTP headers
const EndOfHeaderSeparator = "\r\n\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endOfHeaderIdx := bytes.Index(data, []byte(EndOfHeaderSeparator))
	if endOfHeaderIdx == -1 {
		return 0, false, errors.New("incomplete headers") // incomplete end of headers
	}

	consumedBytes := endOfHeaderIdx + len(EndOfHeaderSeparator)
	requestHeaders := data[:endOfHeaderIdx]
	for header := range bytes.SplitSeq(requestHeaders, []byte("\r\n")) {
		idx := bytes.Index(header, []byte(":"))
		if idx == -1 {
			return 0, false, nil
		}

		fieldName := bytes.TrimSpace(header[:idx])
		if len(fieldName) != len(header[:idx]) {
			return 0, false, errors.New("invalid header")
		}
		fieldValue := bytes.TrimSpace(header[idx+1:])

		fmt.Println("fieldName:", string(fieldName))
		fmt.Println("fieldValue:", string(fieldValue))
		h[string(fieldName)] = string(fieldValue)
	}
	return consumedBytes, true, nil
}
