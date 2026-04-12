package main

import (
	"net/http"
)

func main() {
	api := api{addr: ":8080"}

	mux := http.NewServeMux()

	s := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUserHandler)

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
