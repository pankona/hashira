package functions

import (
	"log"
	"net/http"
	"net/url"

	"github.com/pankona/hashira/hashira-web/functions/hashira"
	"github.com/pankona/hashira/hashira-web/functions/hashira/store"
)

func Call(w http.ResponseWriter, r *http.Request) {
	setHeadersForCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Printf("failed to parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	method, ok := values["method"]
	if !ok {
		log.Printf("method is missing in query string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(method) != 1 {
		log.Printf("only one method is allowed but %d methods are specified", len(method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h := hashira.New(store.NewAccessTokenStore(), store.NewTaskAndPriorityStore())

	switch method[0] {
	case "ping":
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("pong")); err != nil {
			log.Printf("ping failed: %v", err)
		}

	case "test-access-token":
		h.TestAccessToken(w, r)

	case "add":
		h.Add(w, r)

	default:
		log.Printf("%s is not implemented", method[0])
		w.WriteHeader(http.StatusBadRequest)
	}
}

func setHeadersForCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Access-Control-Max-Age", "3600")
}
