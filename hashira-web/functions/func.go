package functions

import (
	"net/http"

	"github.com/pankona/hashira/hashira-web/functions/hashira"
	"github.com/pankona/hashira/hashira-web/functions/hashira/store"
)

func TestAccessToken(w http.ResponseWriter, r *http.Request) {
	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())
	h.TestAccessToken(w, r)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())
	h.Upload(w, r)
}

func Download(w http.ResponseWriter, r *http.Request) {
	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())
	h.Download(w, r)
}
