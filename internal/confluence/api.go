package confluence

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	Client   *http.Client
	Endpoint *url.URL
	user     string
	token    string
}

func NewAPI(email, token, host string) (*API, error) {
	timeout := 30 * time.Second
	return &API{
		Client: &http.Client{
			Timeout: timeout,
		},
		Endpoint: &url.URL{
			Host:   host,
			Path:   "/wiki/rest/api",
			Scheme: "https",
		},
		user:  email,
		token: token,
	}, nil
}

// Build a request and send it to confluence api service.
func (a *API) requestAPI(method, path string, body []byte) (*http.Response, error) {
	switch method {
	case "GET":
		method = http.MethodGet
	case "POST":
		method = http.MethodPost
	case "PUT":
		method = http.MethodPut
	case "DELETE":
		method = http.MethodDelete
	default:
		method = http.MethodGet
	}
	// Add path to URL base path object
	path, _ = url.JoinPath(a.Endpoint.String(), path)
	// Create a Request object
	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		return nil, err
	}
	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Add basic auth to headers
	a.Auth(req)
	// Send request
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
