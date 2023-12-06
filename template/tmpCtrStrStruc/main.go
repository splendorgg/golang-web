package main

import (
	"html/template"
	"net/http"
)

type prodSpec struct {
	Size   string
	Weight float32
	Descr  string
}

type product struct {
	ProdID int
	Name   string
	Cost   float64
	Specs  prodSpec
}

var tpl *template.Template
var prod1 product

func main() {
	prod1 = product{
		ProdID: 15,
		Name:   "Laptop",
		Cost:   900,
		Specs: prodSpec{
			Size:   "150 x 70 x 7 mm",
			Weight: 65,
			Descr:  "Shiny",
		},
	}
	tpl, _ = tpl.ParseGlob("*.html")
	http.HandleFunc("/info", productInfoHandler)
	http.ListenAndServe("", nil)
}

func productInfoHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "productinfo2.html", prod1)
}
