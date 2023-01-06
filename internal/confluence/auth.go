package confluence

import "net/http"

// Auth Adds user and token to request header, implements basic auth.
func (a *API) Auth(req *http.Request) {
	if a.user != "" && a.token != "" {
		req.SetBasicAuth(a.user, a.token)
	}
}
