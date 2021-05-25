package truenasapi

import (
	"fmt"
	"net/http"
)

// AddHeaderTransport holds http Client information
type AddHeaderTransport struct {
	Token string
	T     http.RoundTripper
}

// RoundTrip adds authorization header to http Client
func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", adt.Token))
	return adt.T.RoundTrip(req)
}

func newHttpClient(token string) *http.Client {
	t := http.DefaultTransport
	return &http.Client{Transport: &AddHeaderTransport{
		Token: token,
		T:     t,
	}}
}
