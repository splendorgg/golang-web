package main

import (
	"fmt"
	"net/http"
	"text/template"
)

var tpl *template.Template
var name = "John"

func main() {
	tpl, _ = template.ParseGlob("*.html")
	http.HandleFunc("/welcome", welcomeHandler)
	http.ListenAndServe("", nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("welcomeHandler Running")
	tpl.ExecuteTemplate(w, "welcome.html", name)
}
