package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/pankona/hashira-auth/kvstore"
	"github.com/pankona/hashira-auth/user"
	"github.com/satori/go.uuid"
)

type Twitter struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
	client            *anaconda.TwitterApi
	kvstore           kvstore.KVStore
	callbackURL       string
	credential        map[string]*oauth.Credentials
}

func New(consumerKey, consumerSecret,
	accessToken, accessTokenSecret,
	callbackURL string, kvstore kvstore.KVStore) *Twitter {
	if consumerKey == "" || consumerSecret == "" ||
		accessToken == "" || accessTokenSecret == "" ||
		callbackURL == "" {
		panic("not enough parameter")
	}

	t := &Twitter{
		consumerKey:       consumerKey,
		consumerSecret:    consumerSecret,
		accessToken:       accessToken,
		accessTokenSecret: accessTokenSecret,
		kvstore:           kvstore,
		callbackURL:       callbackURL,
		credential:        make(map[string]*oauth.Credentials),
	}
	t.client = anaconda.NewTwitterApiWithCredentials(
		accessToken, accessTokenSecret,
		consumerKey, consumerSecret)

	return t
}

func (t *Twitter) Register(pattern string) {
	http.Handle(pattern, http.StripPrefix(pattern, t))
}

func (t *Twitter) handleRequestToken(w http.ResponseWriter, r *http.Request) {
	url, tmpCred, err := t.client.AuthorizationURL(t.callbackURL)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	t.credential[tmpCred.Token] = tmpCred
	http.Redirect(w, r, url, http.StatusFound)
}

func (t *Twitter) handleAccessToken(w http.ResponseWriter, r *http.Request) {
	oauthToken := r.URL.Query().Get("oauth_token")
	cred := t.credential[oauthToken]
	c, _, err := t.client.GetCredentials(cred, r.URL.Query().Get("oauth_verifier"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	delete(t.credential, oauthToken)

	cli := anaconda.NewTwitterApiWithCredentials(c.Token, c.Secret, t.consumerKey, t.consumerSecret)

	v := url.Values{}
	v.Set("include_entities", "true")
	v.Set("skip_status", "true")
	u, err := cli.GetSelf(v)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	// check if the user already exists
	uid, ok := t.kvstore.Load("userIDByIDToken", u.IdStr)
	if ok {
		token := uuid.NewV4()
		t.kvstore.Store("userIDByAccessToken", token.String(), uid)
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
		uid, ok = t.kvstore.Load("userIDByAccessToken", a.Value)
		if ok {
			// this user is already registered by other oauth provider
			v, ok := t.kvstore.Load("userByUserID", uid.(string))
			if !ok {
				// fatal
			}
			us := v.(map[string]interface{})
			us["TwitterID"] = u.IdStr
			t.kvstore.Store("userIDByIDToken", u.IdStr, us["ID"])
			t.kvstore.Store("userByUserID", us["ID"].(string), us)
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
	t.kvstore.Store("userIDByIDToken", u.IdStr, userID.String())
	t.kvstore.Store("userByUserID", userID.String(), user.User{
		ID:        userID.String(),
		Name:      username,
		TwitterID: u.IdStr,
	})
	t.kvstore.Store("userIDByAccessToken", token.String(), userID.String())

	cookie := &http.Cookie{
		Name:  "Authorization",
		Value: token.String(),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)

}

func (t *Twitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "callback":
		t.handleAccessToken(w, r)
	default:
		t.handleRequestToken(w, r)
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
