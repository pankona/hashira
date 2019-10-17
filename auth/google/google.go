package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/pankona/hashira/user"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type UserStore interface {
	Store(user *user.User) error
	Fetch(userID string) (*user.User, error)
	FetchByAccessToken(accesstoken string) (*user.User, error)

	FetchByGoogleIDToken(idtoken string) (*user.User, error)
}

// Google is a struct to provide hashira's OIDC functionality using Google
type Google struct {
	id         string
	secret     string
	provider   *oidc.Provider
	verifier   *oidc.IDTokenVerifier
	config     oauth2.Config
	credential map[string]struct{}

	userStore UserStore
}

// New returns Google instance with specified arguments
func New(id, secret, callbackURL string, store UserStore) *Google {
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: id,
	}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  callbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &Google{
		id:         id,
		secret:     secret,
		provider:   provider,
		verifier:   verifier,
		config:     config,
		userStore:  store,
		credential: make(map[string]struct{}),
	}
}

func (g *Google) handleCode(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewV4().String()
	url := g.config.AuthCodeURL(id)
	g.credential[id] = struct{}{}

	http.Redirect(w, r, url, http.StatusFound)
}

func (g *Google) handleIDToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	state := r.URL.Query().Get("state")
	if _, ok := g.credential[state]; !ok {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	// TODO: exclude control
	delete(g.credential, state)

	oauth2Token, err := g.config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	idToken, err := g.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	oauth2Token.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if the user already exists
	us, err := g.userStore.FetchByGoogleIDToken(idToken.Subject)
	if err != nil {
		panic(err)
	}

	if us != nil {
		cookie := &http.Cookie{
			Name:  "Authorization",
			Value: us.AccessToken,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
		return
	}

	// check if the user is registered by other oauth provider
	a, err := r.Cookie("Authorization")
	if err == nil {
		// has Authorization
		us, err := g.userStore.FetchByAccessToken(a.Value)
		if err != nil {
			panic(err)
		}

		if us != nil {
			// this user is already registered by other oauth provider
			// update user to indicate oauth by google has been connected
			us.GoogleID = idToken.Subject
			err = g.userStore.Store(us)
			if err != nil {
				panic(fmt.Sprintf("failed to store user. fatal: %v", err))
			}

			http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
			return
		}
	}

	// create new user
	var (
		userID = uuid.NewV4()
		token  = uuid.NewV4()
	)

	username, err := fetchPhraseFromMashimashi()
	if err != nil {
		panic(fmt.Sprintf("failed to fetch phrase from mashimashi: %v", err))
	}

	err = g.userStore.Store(&user.User{
		ID:          userID.String(),
		Name:        username,
		GoogleID:    idToken.Subject,
		AccessToken: token.String(),
	})
	if err != nil {
		panic(fmt.Errorf("failed store user: %v", err))
	}

	cookie := &http.Cookie{
		Name:  "Authorization",
		Value: token.String(),
		Path:  "/",
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, "http://localhost:3000", http.StatusFound)
}

func (g *Google) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/callback":
		g.handleIDToken(w, r)
	default:
		g.handleCode(w, r)
	}
}

// TODO: make this DRY
func fetchPhraseFromMashimashi() (string, error) {
	resp, err := http.Get("https://strongest-mashimashi.appspot.com/api/v1/phrase")
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
