package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/amrrdev/httpfromtcp/internal/headers"
)

type ParserState string

const (
	StateInit    ParserState = "init"
	StateHeaders ParserState = "headers"
	StateBody    ParserState = "body" // bodys bodys bodys
	StateDone    ParserState = "done"
	StateError   ParserState = "error"
)

type RequestLine struct {
	HttpMethod    string
	RequestTarget string
	HttpVersion   string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        string
	State       ParserState
}

var (
	ErrorMalformRequestLine     = fmt.Errorf("malformed request-line")
	ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
	ErrorRequestInErrorState    = fmt.Errorf("requesr in error state")
	SEPARATOR                   = []byte("\r\n")
)

func GetInt(header *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := header.Get(name)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func NewRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
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
			r.State = StateHeaders
		case StateHeaders:
			bytesRead, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}

			read += bytesRead
			if done {
				length := GetInt(r.Headers, "content-length", 0)
				if length > 0 {
					r.State = StateBody
				} else {
					r.State = StateDone
				}
			}
			if bytesRead == 0 {
				break outer
			}

		case StateBody:
			length := GetInt(r.Headers, "content-length", 0)
			if length == 0 {
				r.State = StateDone
			}

			remaining := min(length-len(r.Body), len(data[read:]))
			if remaining > 0 {
				r.Body += string(data[read : read+remaining])
				read += remaining
			}

			if len(r.Body) >= length {
				r.State = StateDone
			}

			if remaining == 0 {
				break outer
			}
		case StateDone:
			break outer
		default:
			panic("something wendt wrong")
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
		if err != nil && err != io.EOF {
			return nil, err
		}
		bufLen += n

		readN, parseErr := request.Parse(buf[:bufLen])
		if parseErr != nil {
			return nil, parseErr
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

		if err == io.EOF && readN == 0 {
			if request.State == StateBody {
				expectedLength := GetInt(request.Headers, "content-length", 0)
				if len(request.Body) < expectedLength {
					return nil, fmt.Errorf("unexpected EOF: expected body length %d, got %d", expectedLength, len(request.Body))
				}
			}
			break
		}
	}

	return request, nil
}
