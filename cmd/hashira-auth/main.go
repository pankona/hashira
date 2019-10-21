package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pankona/hashira/auth/accesstoken"
	"github.com/pankona/hashira/auth/google"
	"github.com/pankona/hashira/auth/me"
	"github.com/pankona/hashira/auth/store"
	"github.com/pankona/hashira/auth/twitter"
	"gopkg.in/yaml.v2"
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

	buf, err := ioutil.ReadFile("secret.yaml")
	if err != nil {
		panic(err)
	}

	config := map[string]string{}
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	ms := store.NewMemStore()

	googleOAuthHandler := google.New(
		config["GOOGLE_OAUTH2_CLIENT_ID"],
		config["GOOGLE_OAUTH2_CLIENT_SECRET"],
		servingBaseURL+"/auth/google/callback",
		ms)

	twitterOAuthHandler := twitter.New(
		config["TWITTER_API_TOKEN"],
		config["TWITTER_API_SECRET"],
		config["TWITTER_API_ACCESS_TOKEN"],
		config["TWITTER_API_ACCESS_TOKEN_SECRET"],
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

	log.Fatal(http.ListenAndServe(":"+port, r.mux))
}
