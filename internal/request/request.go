package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	reqLine, err := parseRequestLine(req)
	if err != nil {
		return nil, err
	}
	res := Request{
		RequestLine: reqLine,
	}
	return &res, nil
}

func parseRequestLine(req []byte) (RequestLine, error) {
	reqString := string(req)
	reqLine := strings.Split(reqString, "\r\n")[0]
	parts := strings.Split(reqLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("wrong line")
	}
	method, target, version := parts[0], parts[1], parts[2]
	for _, ch := range method {
		if ch < 'A' || ch > 'Z' {
			return RequestLine{}, errors.New("wrong method")
		}
	}
	if version != "HTTP/1.1" {
		return RequestLine{}, errors.New("wrong version")
	}
	res := RequestLine{
		HttpVersion:   "1.1",
		RequestTarget: target,
		Method:        method,
	}
	return res, nil

}
