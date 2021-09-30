package hashira

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var bearerRegex = regexp.MustCompile("(?i)^bearer (.*)$")

func (h *Hashira) TestAccessToken(w http.ResponseWriter, r *http.Request) {
	accesstoken, err := h.retrieveAccessTokenFromHeader(r.Context(), r.Header)
	if err != nil {
		log.Printf("failed to retrieve accesstoken from header: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uid, err := h.AccessTokenStore.FindUidByAccessToken(r.Context(), accesstoken)
	if err != nil {
		log.Printf("could not find a user who has the accesstoken %v: %v", accesstoken, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Printf("accesstoken test for %s (uid: %s) is OK", accesstoken, uid)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Hashira) retrieveAccessTokenFromHeader(ctx context.Context, header http.Header) (string, error) {
	if len(header["Authorization"]) != 1 {
		return "", errors.New("Authorization header is missing or appears more than once. error")
	}

	authHeader := header["Authorization"][0]
	matches := bearerRegex.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", fmt.Errorf("unexpected Authorization header format: %v", authHeader)
	}

	accesstoken := matches[1]
	return accesstoken, nil
}
