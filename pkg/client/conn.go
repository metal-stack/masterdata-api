package client

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Config is the configuration to create a api-server connection
type Config struct {
	BaseURL string
	Token   string
	Debug   bool

	UserAgent string

	// Namespace if set adds this namespace to namespaced requests such that it does not need to be passed all the time
	Namespace string
}

func (d *Config) HttpClient() *http.Client {
	return &http.Client{
		Transport: &AddHeaderTransport{
			debug: d.Debug,
			T:     http.DefaultTransport,
			Token: d.Token,
		},
	}
}

type AddHeaderTransport struct {
	debug bool

	Token string
	T     http.RoundTripper
}

func (a *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+a.Token)

	if a.debug {
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Printf("DEBUG ERROR: %s\n", err)
		} else {
			fmt.Printf("DEBUG REQUEST:\n%s\n", string(reqDump))
		}
	}

	resp, err := a.T.RoundTrip(req)

	if a.debug && resp != nil {
		respDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Printf("DEBUG ERROR: %s\n", err)
		} else {
			fmt.Printf("DEBUG RESPONSE:\n%s\n", string(respDump))
		}
	}

	return resp, err
}
