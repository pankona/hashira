package hashira

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Hashira) Download(w http.ResponseWriter, r *http.Request) {
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

	tp, err := h.TaskAndPriorityStore.Load(r.Context(), uid)
	if err != nil {
		log.Printf("failed to load tasks and priorities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tp); err != nil {
		log.Printf("failed to write response body: %v", err)
		return
	}
}
