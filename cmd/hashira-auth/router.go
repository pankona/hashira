package main

import (
	"fmt"
	"log"
	"net/http"
)

type router struct {
	mux    *http.ServeMux
	routes map[http.Handler]*routes
}

type routes struct {
	path      string
	subroutes []string
}

func (r *router) route() {
	for handler, routes := range r.routes {
		for _, subroute := range routes.subroutes {
			r.mux.Handle(fmt.Sprintf("%s%s", routes.path, subroute),
				http.StripPrefix(routes.path, handler))
		}
	}
}
