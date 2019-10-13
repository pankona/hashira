package accesstoken

import "net/http"

// AccessToken is a struct that implements http.Handler for resource "accesstoken"
type AccessToken struct{}

// New returns a struct that implements http.handler for resource "accesstoken"
func New() *AccessToken {
	return &AccessToken{}
}

func (a *AccessToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
	// POST request to generate new access token
}
