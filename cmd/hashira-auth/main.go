package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/pankona/hashira/auth/google"
	"github.com/pankona/hashira/auth/twitter"
	"github.com/pankona/hashira/kvstore"
	"github.com/pankona/hashira/user"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	env := os.Getenv("GAE_ENV")
	servingBaseURL := "http://localhost:8080"
	if env != "" {
		servingBaseURL = "https://hashira-auth.appspot.com"
	}

	log.Printf("GAE_ENV: %v", env)
	log.Printf("servingBaseURL: %v", servingBaseURL)

	kvs := &kvstore.DSStore{}
	registerGoogle(kvs, servingBaseURL)
	registerTwitter(kvs, servingBaseURL)

	http.HandleFunc("/api/v1/accesstoken", handleAccessToken)
	http.HandleFunc("/api/v1/me", handleMe)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleAccessToken(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	// POST request to generate new access token
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "origin,x-requested-with,content-type,accept,authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	kvs := &kvstore.DSStore{}

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		// No Authorization info found
		w.WriteHeader(404)
		_, err := w.Write([]byte("no authorization info found"))
		if err != nil {
			log.Printf("failed to write response: %v", err)
		}
		return
	}

	userID, ok := kvs.Load("userIDByAccessToken", auth)
	if !ok {
		// UserID that has specified access token not found
		w.WriteHeader(404)
		return
	}

	u, ok := kvs.Load("userByUserID", userID.(string))
	if !ok {
		// User that has specified User ID not found
		w.WriteHeader(404)
		return
	}

	m, ok := u.(map[string]interface{})
	if !ok {
		w.WriteHeader(500)
		return
	}
	me := user.User{
		ID:   m["ID"].(string),
		Name: m["Name"].(string),
		GoogleID: func() string {
			if m["GoogleID"] != "" {
				return "***"
			}
			return ""
		}(),
		TwitterID: func() string {
			if m["GoogleID"] != "" {
				return "***"
			}
			return ""
		}(),
	}

	buf, err := json.Marshal(me)
	if err != nil {
		// Internal server error
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		// failed write response. Just logging
		log.Printf("failed to write response: %v", err)
	}

	w.WriteHeader(200)
}

func registerGoogle(kvs kvstore.KVStore, servingBaseURL string) {
	var (
		clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
		clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
	)
	g := google.New(clientID, clientSecret,
		servingBaseURL+"/auth/google/callback", kvs)
	g.Register("/auth/google/")
}

func registerTwitter(kvs kvstore.KVStore, servingBaseURL string) {
	var (
		consumerKey       = os.Getenv("TWITTER_API_TOKEN")
		consumerSecret    = os.Getenv("TWITTER_API_SECRET")
		accessToken       = os.Getenv("TWITTER_API_ACCESS_TOKEN")
		accessTokenSecret = os.Getenv("TWITTER_API_ACCESS_TOKEN_SECRET")
	)
	t := twitter.New(consumerKey, consumerSecret, accessToken, accessTokenSecret,
		servingBaseURL+"/auth/twitter/callback", kvs)
	t.Register("/auth/twitter/")
}
