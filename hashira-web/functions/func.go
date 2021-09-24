package functions

import (
	"net/http"

	"github.com/pankona/hashira/hashira-web/functions/hashira"
)

func TestAccessToken(w http.ResponseWriter, r *http.Request) {
	h := hashira.New()
	h.TestAccessToken(w, r)
}
