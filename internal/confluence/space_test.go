package confluence

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetSpace(t *testing.T) {
	testCases := []struct {
		desc string
		key  string
	}{
		{
			desc: "Get Space request",
			key:  "devops",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			server := httpServer([]byte(`{
				"id":65541,"key":"DEVOPS","name":"devops","type":"global","status":"current"
			}`))
			defer server.Close()
			serverURL := strings.Split(server.URL, "://")

			api := &API{
				Client: server.Client(),
				Endpoint: &url.URL{
					Scheme: serverURL[0],
					Host:   serverURL[1],
				},
			}

			space, err := api.GetSpace("devops")
			if err != nil {
				t.Error(err)
			}
			if space == nil {
				t.Error("Space struct got empty")
			}
			if space.Id != 65541 {
				t.Errorf("wants space.id: %v, but got %v", 65541, space.Id)
			}
			if space.Key != "DEVOPS" {
				t.Errorf("wants space.key: %v, but got %v", "DEVOPS", space.Key)
			}
		})
	}
}

func httpServer(resp []byte) *httptest.Server {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))
	return server
}
