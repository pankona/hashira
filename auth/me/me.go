// Package me is a package that has responsibility for resource me
package me

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pankona/hashira/user"
)

type UserStore interface {
	FetchByAccessToken(accesstoken string) (*user.User, error)
}

// Me is a struct that implements http.Handler for resource "me"
type Me struct {
	userStore UserStore
}

// New returns a struct that implements http.handler for resource "me"
func New(s UserStore) *Me {
	return &Me{
		userStore: s,
	}
}

func (m *Me) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "origin,x-requested-with,content-type,accept,authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		// No Authorization info found
		w.WriteHeader(404)
		_, err := w.Write([]byte("no authorization info found"))
		if err != nil {
			log.Printf("failed to write response: %v", err)
			w.WriteHeader(500)
		}
		return
	}

	u, err := m.userStore.FetchByAccessToken(auth)
	if err != nil {
		panic(err)
	}

	if u == nil {
		// user that has specified access token not found
		w.WriteHeader(404)
		return
	}

	buf, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(buf)
	if err != nil {
		// failed write response. Just logging
		log.Printf("failed to write response: %v", err)
		w.WriteHeader(500)
		return
	}
}
