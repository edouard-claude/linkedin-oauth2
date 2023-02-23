package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/oauth2"
)

func init() {
	os.Setenv("LINKEDIN_CLIENT_ID", "mock-client-id")
	os.Setenv("LINKEDIN_CLIENT_SECRET", "mock-client-secret")
	os.Setenv("REDIRECT_URL", "http://localhost:8080/callback")
}

// Test the displayProfile function
func TestDisplayProfile(t *testing.T) {
	// Create a mock HTTP server to return a sample response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/me" {
			profile := map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"email":     "johndoe@example.com",
			}
			json.NewEncoder(w).Encode(profile)
		}
	}))
	defer server.Close()

	// Create a new OAuth2 token with a mock access token
	token := &oauth2.Token{AccessToken: "mock-access-token"}

	// Update the profileEndpoint to use the mock server URL
	profileEndpoint = server.URL + "/v2/me"

	// Call the displayProfile function with the mock token
	profile := displayProfile(token)

	// Verify that the profile was parsed correctly
	if profile == nil {
		t.Error("Failed to parse profile information")
	}
	if profile["firstName"] != "John" {
		t.Errorf("Expected firstName to be 'John', but got '%s'", profile["firstName"])
	}
	if profile["lastName"] != "Doe" {
		t.Errorf("Expected lastName to be 'Doe', but got '%s'", profile["lastName"])
	}
	if profile["email"] != "johndoe@example.com" {
		t.Errorf("Expected email to be 'johndoe@example.com', but got '%s'", profile["email"])
	}
}

func TestIndexHandler(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/", nil)

	// Create a new HTTP response recorder to capture the response
	w := httptest.NewRecorder()

	// Call the indexHandler function with the mock request and response
	indexHandler(w, req)

	// Verify that the response is a redirect to the OAuth2 provider
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status code %d, but got %d", http.StatusTemporaryRedirect, w.Code)
	}
	expectedLocation := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	if w.Header().Get("Location") != expectedLocation {
		t.Errorf("Expected Location header to be '%s', but got '%s'", expectedLocation, w.Header().Get("Location"))
	}
}
