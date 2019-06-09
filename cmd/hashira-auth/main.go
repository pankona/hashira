package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pankona/hashira-auth/google"
	"github.com/pankona/hashira-auth/kvstore"
	"github.com/pankona/hashira-auth/twitter"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		msg := "<html>"
		msg += "<head><title>hashira-auth</title></head>"
		msg += "<body>"

		a, err := r.Cookie("Authorization")
		if err != nil {
			msg += fmt.Sprintf("No Authorization info found...<br>")
			msg += fmt.Sprintf("Cookies: %v<br>", r.Cookies())
			msg += "<a href=/auth/google>login by google</a><br>"
			msg += "<a href=/auth/twitter>login by twitter</a><br>"
			msg += "</html>"
			w.Write([]byte(msg))
			return
		}

		userID, ok := kvs.Load("userIDByAccessToken", a.Value)
		if !ok {
			msg += fmt.Sprintf("UserID that has access token [%s] not found...<br>", a.Value)
			msg += "<a href=/auth/google>login by google</a><br>"
			msg += "<a href=/auth/twitter>login by twitter</a><br>"
			msg += "</html>"
			w.Write([]byte(msg))
			return
		}

		u, ok := kvs.Load("userByUserID", userID.(string))
		if !ok {
			msg += fmt.Sprintf("User that has user ID [%s] not found...<br>", userID.(string))
			msg += "<a href=/auth/google>login by google</a><br>"
			msg += "<a href=/auth/twitter>login by twitter</a><br>"
			msg += "</html>"
			w.Write([]byte(msg))
			return
		}
		msg += fmt.Sprintf("Hello, %s!<br>", u.(map[string]interface{})["Name"])

		msg += "<a href=/auth/google>login by google</a>"
		if u.(map[string]interface{})["GoogleID"] != "" {
			msg += " Connected!<br>"
		} else {
			msg += "<br>"
		}

		msg += "<a href=/auth/twitter>login by twitter</a>"
		if u.(map[string]interface{})["TwitterID"] != "" {
			msg += " Connected!<br>"
		} else {
			msg += "<br>"
		}

		msg += "</html>"
		w.Write([]byte(msg))
	})

	http.HandleFunc("/api/v1/accesstoken", func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement
		// POST request to generate new access token
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
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
