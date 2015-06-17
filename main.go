package dumper

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type PrintType int8

const (
	Request PrintType = iota
	Response
)

type Transport struct {
	http.RoundTripper
	Print func(PrintType, []byte)
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.RoundTripper
	if t.Print == nil {
		return rt.RoundTrip(req)
	}

	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	t.Print(Request, b)

	res, err := rt.RoundTrip(req)
	if err != nil {
		return res, err
	}

	b, err = httputil.DumpResponse(res, true)
	if err != nil {
		return res, err
	}

	t.Print(Response, b)

	return res, nil

}

func Print(t PrintType, b []byte) {
	switch t {
	case Request:
		fmt.Printf("REQUEST\n%s\n\n", b)
	case Response:
		fmt.Printf("RESPONSE\n%s\n\n", b)
	}
}
