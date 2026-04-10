package main

import (
	"net/http"
)

type api struct {
	addr string
}

func (a api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sent a list of users"))
}

func (a api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Created user"))
}

func main() {
	api := api{addr: ":8080"}

	mux := http.NewServeMux()

	s := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUserHandler)

	s.ListenAndServe()
}
