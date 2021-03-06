package twitter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/pankona/hashira/auth/user"

	uuid "github.com/satori/go.uuid"
)

type UserStore interface {
	Store(user *user.User) error
	Fetch(userID string) (*user.User, error)
	FetchByAccessToken(accesstoken string) (*user.User, error)

	FetchByTwitterIDToken(idtoken string) (*user.User, error)
}

// Twitter is a struct to provide hashira's oauth functionality using twitter
type Twitter struct {
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
	client            *anaconda.TwitterApi
	callbackURL       string
	credential        map[string]*oauth.Credentials

	userStore UserStore
}

// New returns Twitter instance with specified arguments
func New(consumerKey, consumerSecret, accessToken, accessTokenSecret,
	callbackURL string, store UserStore) *Twitter {
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
		userStore:         store,
		callbackURL:       callbackURL,
		credential:        make(map[string]*oauth.Credentials),
	}

	t.client = anaconda.NewTwitterApiWithCredentials(
		accessToken, accessTokenSecret,
		consumerKey, consumerSecret)

	return t
}

// Register registers an endpoint for Twitter oauth
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
	us, err := t.userStore.FetchByTwitterIDToken(u.IdStr)
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
		http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
		return
	}

	// check if the user is registered by other oauth provider
	a, err := r.Cookie("Authorization")
	if err == nil {
		// has Authorization
		us, err := t.userStore.FetchByAccessToken(a.Value)
		if err != nil {
			panic(err)
		}

		if us != nil {
			// this user is already registered by other oauth provider
			// update user to indicate oauth by twitter has been connected
			us.TwitterID = u.IdStr
			err = t.userStore.Store(us)
			if err != nil {
				panic(fmt.Sprintf("failed to store user. fatal: %v", err))
			}

			http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
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

	err = t.userStore.Store(&user.User{
		ID:          userID.String(),
		Name:        username,
		TwitterID:   u.IdStr,
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
	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
}

func (t *Twitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/callback":
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
