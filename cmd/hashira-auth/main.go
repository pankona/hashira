package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pankona/hashira/kvstore"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	env := os.Getenv("GAE_ENV")
	servingBaseURL := "http://localhost:8080"
	if env != "" {
		servingBaseURL = "https://hashira-auth.appspot.com"
	}

	log.Printf("GAE_ENV: %v", env)
	log.Printf("servingBaseURL: %v", servingBaseURL)

	r := router{
		mux:            http.DefaultServeMux,
		kvs:            &kvstore.DSStore{},
		servingBaseURL: servingBaseURL,
	}
	r.route()

	log.Fatal(http.ListenAndServe(":"+port, r.mux))
}
