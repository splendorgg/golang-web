package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	Name   string
	Image  string
	Price  int
	Rating int
}
type Charts struct {
	Type string
	Data string
}

func main() {
	//Mysql connection

	router := http.NewServeMux()
	router.HandleFunc("/api/column", Columnhandler)
	router.HandleFunc("/api/donut", Donuthandler)
	router.HandleFunc("/api/apache", Apachehandler)
	router.HandleFunc("/api/productcard", Productcard)
	http.ListenAndServe(":3030", router)

}

func Columnhandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	dsn := "root:PASSWORD@tcp(127.0.0.1:3306)/svelte_dashboard"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var data Charts
	db.Where("type = ?", "column").First(&data)
	fmt.Fprintf(w, data.Data)

}

func Donuthandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	dsn := "root:PASSWORD@tcp(127.0.0.1:3306)/svelte_dashboard"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var data Charts
	db.Where("type = ?", "donut").First(&data)
	fmt.Fprintf(w, data.Data)

}
func Apachehandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	dsn := "root:PASSWORD@tcp(127.0.0.1:3306)/svelte_dashboard"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var data Charts
	db.Where("type = ?", "apache").First(&data)
	fmt.Fprintf(w, data.Data)

}

func setHeaders(w *http.ResponseWriter) { // middleware

	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")

}

func Productcard(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w)
	dsn := "root:PASSWORD@tcp(127.0.0.1:3306)/svelte_dashboard"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var products []Product
	db.Find(&products)
	json.NewEncoder(w).Encode(products)

}
