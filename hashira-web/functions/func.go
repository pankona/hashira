package functions

import (
	"net/http"

	"github.com/pankona/hashira/hashira-web/functions/hashira"
	"github.com/pankona/hashira/hashira-web/functions/hashira/store"
)

func TestAccessToken(w http.ResponseWriter, r *http.Request) {
	setHeadersForCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())
	h.TestAccessToken(w, r)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	setHeadersForCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())
	h.Upload(w, r)
}

func Download(w http.ResponseWriter, r *http.Request) {
	setHeadersForCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())
	h.Download(w, r)
}

func setHeadersForCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")
}
