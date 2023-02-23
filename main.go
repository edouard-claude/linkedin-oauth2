package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

// Configure the OAuth2 client
var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
	ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
	Scopes:       []string{"r_liteprofile", "r_emailaddress", "w_member_social"},
	Endpoint:     linkedin.Endpoint,
	RedirectURL:  os.Getenv("REDIRECT_URL"),
}

var profileEndpoint = "https://api.linkedin.com/v2/me"

func displayProfile(token *oauth2.Token) map[string]interface{} {
	// Query the LinkedIn API for the user's profile information
	req, _ := http.NewRequest("GET", profileEndpoint, nil)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to fetch profile information:", err)
		return nil
	}
	defer res.Body.Close()

	// Parse the response and extract the profile data
	var profile map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		fmt.Println("Failed to parse profile information:", err)
		return nil
	}

	return profile
}

// Serve the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Handle the callback from the OAuth2 provider
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	profile := displayProfile(token)
	out, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		fmt.Println("Failed to parse profile information:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func main() {
	// get if env variables are set
	if os.Getenv("LINKEDIN_CLIENT_ID") == "" {
		fmt.Println("LINKEDIN_CLIENT_ID is not set")
		os.Exit(1)
	}
	if os.Getenv("LINKEDIN_CLIENT_SECRET") == "" {
		fmt.Println("LINKEDIN_CLIENT_SECRET is not set")
		os.Exit(1)
	}
	if os.Getenv("REDIRECT_URL") == "" {
		fmt.Println("REDIRECT_URL is not set")
		os.Exit(1)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/callback", callbackHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
