package trhttp

import "net/http"

type trHTTP struct {
	addr string
	cl   *http.Client
}

func New(addr string) *trHTTP {
	cl := &http.Client{}
	return &trHTTP{
		addr: addr,
		cl:   cl,
	}
}
