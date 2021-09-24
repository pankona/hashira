package hashira

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/pankona/hashira/hashira-web/functions/hashira/store"
)

var bearerRegex = regexp.MustCompile("(?i)^bearer (.*)$")

type AccessTokenStore interface {
	FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error)
}

type Hashira struct {
	AccessTokenStore AccessTokenStore
}

func New() *Hashira {
	return &Hashira{
		AccessTokenStore: store.NewAccessTokenStore(),
	}
}

func (h *Hashira) TestAccessToken(w http.ResponseWriter, r *http.Request) {
	if len(r.Header["Authorization"]) != 1 {
		log.Println("Authorization header is missing or appears more than once. error.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authHeader := r.Header["Authorization"][0]
	matches := bearerRegex.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		log.Printf("unexpected Authorization header format: %v", authHeader)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accesstoken := matches[1]
	uid, err := h.AccessTokenStore.FindUidByAccessToken(r.Context(), accesstoken)
	if err != nil {
		log.Printf("could not find a user who has the accesstoken %v: %v", authHeader, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("accesstoken test for %s (uid: %s) is OK", accesstoken, uid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
