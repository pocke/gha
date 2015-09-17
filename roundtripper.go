package gha

import "net/http"

type RoundTripper string

func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+string(rt))
	return http.DefaultTransport.RoundTrip(req)
}

var _ http.RoundTripper = RoundTripper("")
