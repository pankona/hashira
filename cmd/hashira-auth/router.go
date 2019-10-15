package main

import (
	"net/http"
	"os"

	"github.com/pankona/hashira/auth/accesstoken"
	"github.com/pankona/hashira/auth/google"
	"github.com/pankona/hashira/auth/me"
	"github.com/pankona/hashira/auth/twitter"
	"github.com/pankona/hashira/store"
)

type router struct {
	mux            *http.ServeMux
	store          store.Store
	servingBaseURL string
}

func (r *router) route() {
	r.handleMe("/api/v1/me")
	r.handleAccessToken("/api/v1/accesstoken")
	r.handleGoogleOAuth(r.servingBaseURL, "/auth/google")
	r.handleTwitterOAuth(r.servingBaseURL, "/auth/twitter")
}

func (r *router) handleMe(path string) {
	r.mux.Handle(
		path,
		http.StripPrefix(path, me.New(r.store)))
}

func (r *router) handleAccessToken(path string) {
	r.mux.Handle(
		path,
		http.StripPrefix(path, accesstoken.New()))
}

func (r *router) handleGoogleOAuth(servingBaseURL, path string) {
	g := google.New(
		os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"),
		os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"),
		servingBaseURL+path+"/callback",
		r.store)

	r.mux.Handle(
		path,
		http.StripPrefix(path, g))

	r.mux.Handle(
		path+"/callback",
		http.StripPrefix(path, g))
}

func (r *router) handleTwitterOAuth(servingBaseURL, path string) {
	t := twitter.New(
		os.Getenv("TWITTER_API_TOKEN"),
		os.Getenv("TWITTER_API_SECRET"),
		os.Getenv("TWITTER_API_ACCESS_TOKEN"),
		os.Getenv("TWITTER_API_ACCESS_TOKEN_SECRET"),
		servingBaseURL+path+"/callback",
		r.store)

	r.mux.Handle(
		path,
		http.StripPrefix(path, t))

	r.mux.Handle(
		path+"/callback",
		http.StripPrefix(path, t))
}
