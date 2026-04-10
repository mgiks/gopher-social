package main

import (
	"log"
	"net/http"
)

type server struct {
	addr string
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("index page\n"))
			return
		case "/users":
			w.Write([]byte("users page\n"))
			return
		}
	default:
		w.Write([]byte("404 page\n"))
		return
	}
}

func main() {
	s := server{addr: ":8080"}
	log.Fatal(http.ListenAndServe(s.addr, s))
}
