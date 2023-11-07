package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func main() {
	// func ParseFiles(filenames ...string) (*Template, error)
	tpl, _ = template.ParseFiles("index1.html")
	// func (t *Template) ParseFiles(filenames ...string) (*Template, error)
	// tpl, _ = tpl.ParseFiles("index1.html")
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// func (t *Template) Execute(wr io.Writer, data interface{}) error
	tpl.Execute(w, nil)
}
