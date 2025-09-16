package request

import (
	"bytes"
	"fmt"
	"io"
)

type ParserState string

const (
	StateInit  ParserState = "init"
	StateDone  ParserState = "done"
	StateError ParserState = "error"
)

type RequestLine struct {
	HttpMethod    string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

var (
	ErrorMalformRequestLine     = fmt.Errorf("malformed request-line")
	ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
	ErrorRequestInErrorState    = fmt.Errorf("requesr in error state")
	SEPARATOR                   = []byte("\r\n")
)

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r *Request) Done() bool {
	return r.State == StateDone
}

func (r *Request) Error() bool {
	return r.State == StateError
}

func (r *Request) Parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.State {
		case StateError:
			return 0, ErrorRequestInErrorState

		case StateInit:
			rl, n, err := ParseRequestLine(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = StateDone
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func ParseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	consumedBytes := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorMalformRequestLine
	}

	rl := &RequestLine{
		HttpMethod:    string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, consumedBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.Done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.Parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil
}
