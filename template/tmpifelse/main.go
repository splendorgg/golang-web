package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

// User first letter must be capitalized to be exported
type User struct {
	Name     string
	Language string
	Member   bool
}

var u User

func main() {
	//u = User{Name: "Abuzer", Language: "English", Member: false}
	//u = User{Name: "Buzer", Language: "Spanish", Member: true}
	u = User{Name: "Uzer", Language: "", Member: true}

	tpl, _ = tpl.ParseGlob("*.html")
	http.HandleFunc("/welcome", welcomeHandler)
	http.ListenAndServe("", nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "membership2.html", u)
}
