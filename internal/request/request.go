package request

import (
	"errors"
	"io"
	"strings"
)

type RequestStatus int

const (
	Initialized = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	Status      RequestStatus
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req Request
	req.Status = Initialized
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	for req.Status != Done {
		if readToIndex == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		numBytes, err := reader.Read(buffer[readToIndex:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		if numBytes == 0 && err == io.EOF {
			req.Status = Done
			break
		}
		readToIndex += numBytes
		n, err := req.parse(buffer)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			copy(buffer, buffer[n:readToIndex])
			readToIndex -= n
		}
	}
	return &req, nil
}

func parseRequestLine(data []byte) (method, target, version string, consumed int, err error) {
	// Find the end of the request line
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			// Found the end of request line
			line := string(data[:i])
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return "", "", "", 0, errors.New("malformed request line")
			}
			return parts[0], parts[1], parts[2], i + 2, nil // +2 for \r\n
		}
	}
	// No complete request line found
	return "", "", "", 0, nil
}

// parse parses the next chunk of data
func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case Initialized:
		method, target, version, consumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if consumed == 0 {
			// Need more data
			return 0, nil
		}

		r.RequestLine = RequestLine{
			Method:        method,
			RequestTarget: target,
			HttpVersion:   version[5:],
		}
		r.Status = Done
		return consumed, nil

	case Done:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}
