package hashira

import (
	"encoding/json"
	"log"
	"net/http"
)

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
