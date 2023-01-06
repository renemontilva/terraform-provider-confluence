package confluence

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestAPI(t *testing.T) {
	JsonResponse := []byte(`
	{
  		"id": "12345",
  		"title": "Terrafrom Test",
  		"type": "page",
  		"space": {
  		  "id": 65541,
		  "key": "DEVOPS",
  		  "name": "devops",
		  "status": "current"
		}
	}
	`)

	JsonRequest := []byte(`
	{
  		"title": "Terrafrom Test",
  		"type": "page",
  		"space": {
		  "key": "DEVOPS",
		}
	}
	`)
	// Init a test server with an expected response
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(JsonResponse)
	}))
	defer server.Close()
	// Test RequestAPI calls
	serverURL := strings.Split(server.URL, "://")
	client := &API{
		Client: server.Client(),
		Endpoint: &url.URL{
			Scheme: serverURL[0],
			Host:   serverURL[1],
		},
		user:  "user@email.com",
		token: "123456",
	}

	resp, err := client.requestAPI(http.MethodPost, "/content", JsonRequest)
	if err != nil {
		t.Errorf("requestAPI error: %v", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ReadAll error: %v", err)
	}
	var content Content
	err = json.Unmarshal(b, &content)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	if content.Id != "12345" {
		t.Errorf("Content Id, wants 12345, but got %v", content.Id)
	}

	if content.Space.Key != "DEVOPS" {
		t.Errorf("Content space key, wants DEVOPS, but got %v", content.Space.Key)
	}
}
