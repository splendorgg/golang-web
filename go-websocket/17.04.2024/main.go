package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Id   int      `json:"id" `
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

var user = User{
	Id:   1,
	Name: "Abuzer",
	Tags: []string{"qwe", "asd", "123"},
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/json", jsonHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.ExecuteTemplate(w, "index.html", user)

}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("json.html")
	data, _ := json.Marshal(user) 
	tmpl.ExecuteTemplate(w, "json.html", template.JS(data))


}


