package google

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/pankona/hashira-auth/kvstore"
	"github.com/pankona/hashira-auth/user"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type Google struct {
	id       string
	secret   string
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	config   oauth2.Config
	kvstore  kvstore.KVStore
}

func New(id, secret, callbackURL string, kvstore kvstore.KVStore) *Google {
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
		id:       id,
		secret:   secret,
		provider: provider,
		verifier: verifier,
		config:   config,
		kvstore:  kvstore,
	}
}

func (g *Google) Register(pattern string) {
	http.Handle(pattern, http.StripPrefix(pattern, g))
}

var state = "foobar" // Don't do this in production.

func (g *Google) handleCode(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, g.config.AuthCodeURL(state), http.StatusFound)
}

func (g *Google) handleIDToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.URL.Query().Get("state") != state {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

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
	uid, ok := g.kvstore.Load("userIDByIDToken", idToken.Subject)
	if ok {
		token := uuid.NewV4()
		g.kvstore.Store("userIDByAccessToken", token.String(), uid)
		cookie := &http.Cookie{
			Name:  "Authorization",
			Value: token.String(),
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// check if the user is registered by other oauth provider
	a, err := r.Cookie("Authorization")
	if err == nil {
		// has Authorization
		uid, ok = g.kvstore.Load("userIDByAccessToken", a.Value)
		if ok {
			// this user is already registered by other oauth provider
			v, ok := g.kvstore.Load("userByUserID", uid.(string))
			if !ok {
				// fatal
			}
			us := v.(user.User)
			us.GoogleID = idToken.Subject
			g.kvstore.Store("userIDByIDToken", idToken.Subject, us.ID)
			g.kvstore.Store("userByUserID", us.ID, us)
			http.Redirect(w, r, "/", http.StatusFound)
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
		// TODO: error handling
	}

	g.kvstore.Store("userIDByIDToken", idToken.Subject, userID.String())
	g.kvstore.Store("userByUserID", userID.String(), user.User{
		ID:       userID.String(),
		Name:     username,
		GoogleID: idToken.Subject,
	})
	g.kvstore.Store("userIDByAccessToken", token.String(), userID.String())

	cookie := &http.Cookie{
		Name:  "Authorization",
		Value: token.String(),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (g *Google) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "callback":
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
		if resp != nil {
			resp.Body.Close()
		}
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
