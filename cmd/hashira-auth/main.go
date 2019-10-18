package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pankona/hashira/auth/accesstoken"
	"github.com/pankona/hashira/auth/google"
	"github.com/pankona/hashira/auth/me"
	"github.com/pankona/hashira/auth/twitter"
	"github.com/pankona/hashira/store"
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

	ms := store.NewMemStore()

	googleOAuthHandler := google.New(
		os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"),
		os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"),
		servingBaseURL+"/auth/google/callback",
		ms)

	twitterOAuthHandler := twitter.New(
		os.Getenv("TWITTER_API_TOKEN"),
		os.Getenv("TWITTER_API_SECRET"),
		os.Getenv("TWITTER_API_ACCESS_TOKEN"),
		os.Getenv("TWITTER_API_ACCESS_TOKEN_SECRET"),
		servingBaseURL+"/auth/twitter/callback",
		ms)

	r := router{
		mux: http.DefaultServeMux,
		routes: map[http.Handler]*routes{
			me.New(ms): {
				"/api/v1/me",
				[]string{
					"",
				},
			},
			accesstoken.New(): {
				"/api/v1/accesstoken",
				[]string{
					"",
				},
			},
			googleOAuthHandler: {
				"/auth/google",
				[]string{
					"",
					"/callback",
				},
			},
			twitterOAuthHandler: {
				"/auth/twitter",
				[]string{
					"",
					"/callback",
				},
			},
		},
	}

	r.route()
	//r.mux.Handle("/", &root{})

	log.Fatal(http.ListenAndServe(":"+port, r.mux))
}

type root struct{}

func (ro *root) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("r.URL.Path: %v", r.URL.Path)
}
