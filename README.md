# HTTPFromTCP

> Building HTTP from the ground up, one byte at a time.

A Go implementation that parses HTTP requests directly from raw TCP connections.

## Overview

HTTP parser built from scratch that reads TCP byte streams and constructs HTTP request objects. Implements RFC 7230 compliant parsing for request lines, headers, and message structure.

```bash
curl localhost:42069/api/users
```

Raw TCP stream:

```http
GET /api/users HTTP/1.1\r\n
Host: localhost:42069\r\n
Accept: application/json\r\n
\r\n
```

Parsed output:

```go
Request{
  Method: "GET",
  Path: "/api/users",
  Version: "1.1",
  Headers: map[string]string{
    "host": "localhost:42069",
    "accept": "application/json",
  }
}
```

## Current Status

### Work in Progress

**What works:**

- [x] TCP listener on port 42069
- [x] HTTP request line parsing (method, path, version)
- [x] HTTP headers parsing and validation
- [x] Basic request structure

**Coming soon:**

- [ ] Request body parsing
- [ ] HTTP response generation
- [ ] Chunked transfer encoding

## Quick Start

```bash
# Clone and run the TCP listener
git clone https://github.com/amrrdev/httpfromtcp
cd httpfromtcp
go run cmd/tcplistener/main.go

# In another terminal, send a request
curl localhost:42069/hello
```

## Why?

Because understanding the fundamentals is cool. And sometimes you need to know exactly what's happening under the hood.

## Contributing

This is a learning project, but PRs are welcome! Especially if you want to help with:

- Chunked encoding support
- Response generation
- More robust error handling

---

Made with ☕ and curiosity
