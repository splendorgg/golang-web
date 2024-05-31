package main

import (
	"html/template"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		tpl, _ := template.ParseFiles("index.html")
		tpl.Execute(w, struct{ UserID string }{UserID: userID})
		w.Write([]byte("User ID: " + userID))
	})
	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	log.Printf("Server has started %s", s.addr)
	return server.ListenAndServe()

}
