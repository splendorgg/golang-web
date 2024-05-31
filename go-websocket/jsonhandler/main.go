package main

import (
	"encoding/json"
	"net/http"
)

type user struct {
	Id   int
	Name string
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /api", handler)
	http.ListenAndServe(":3030", router)

}

func handler(w http.ResponseWriter, r *http.Request) {
	u := user{
		Id:   1,
		Name: "Abuzer",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}
