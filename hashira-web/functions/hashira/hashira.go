package hashira

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var bearerRegex = regexp.MustCompile("(?i)^bearer (.*)$")

type AccessTokenStore interface {
	FindUidByAccessToken(ctx context.Context, accesstoken string) (string, error)
}

type TaskAndPriorityStore interface {
	Save(ctx context.Context, uid string, tp TaskAndPriority) error
}

type Hashira struct {
	AccessTokenStore     AccessTokenStore
	TaskAndPriorityStore TaskAndPriorityStore
}

func New(atStore AccessTokenStore, tpStore TaskAndPriorityStore) *Hashira {
	return &Hashira{
		AccessTokenStore:     atStore,
		TaskAndPriorityStore: tpStore,
	}
}

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

type Task struct {
	ID        string
	Name      string
	Place     string
	IsDeleted bool
}

type Priority map[string][]string

// priority's key should be one of following strings:
// "BACKLOG", "TODO", "DOING", "DONE"
type TaskAndPriority struct {
	Tasks    map[string]Task `json:"tasks"`
	Priority Priority        `json:"priority"`
}

func (h *Hashira) Upload(w http.ResponseWriter, r *http.Request) {
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

	var tp TaskAndPriority
	err = json.NewDecoder(r.Body).Decode(&tp)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.TaskAndPriorityStore.Save(r.Context(), uid, tp); err != nil {
		log.Printf("failed to save tasks and priorities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
