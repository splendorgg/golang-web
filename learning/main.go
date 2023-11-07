package main

import (
	"fmt"
	"net/http"
	"time"
)

func helloWorldPage(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "Hello World")
	case "/deneme":
		fmt.Fprint(w, "Deneme")
	default:
		fmt.Fprint(w, "ERROR! ERROR!")
		fmt.Println("ERROR! ERROR!")
	}
	fmt.Printf("Handling function with %s request \n", r.Method)
}

func htmlVsPlain(w http.ResponseWriter, r *http.Request) {
	fmt.Println("htmlVsPlain")
	//w.Header().Set("Content-Type", "text/plain") // <h1>Hello World</h1>
	w.Header().Set("Content-Type", "text/html") // Hello World
	fmt.Fprint(w, "<h1>Hello World</h1>")

}
func timeout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Timeout attempt")
	time.Sleep(2 * time.Second)
	fmt.Fprint(w, "Did *not* timeout")

}
func helloWorldSpecial(w http.ResponseWriter, r *http.Request) {
	fmt.Println("helloWorldSpecial")
	fmt.Fprint(w, "<h1 style=\"background-color:grey;\">Hello World</h1>")

}

func main() {
	http.HandleFunc("/", helloWorldPage)
	http.HandleFunc("/html", htmlVsPlain)
	http.HandleFunc("/timeout", timeout)

	server := http.Server{
		Addr:         "",
		Handler:      nil,
		ReadTimeout:  1000,
		WriteTimeout: 1000,
	}

	var specialMux http.ServeMux
	server.Handler = &specialMux
	specialMux.HandleFunc("/special", helloWorldSpecial)
	//http.ListenAndServe("", nil)
	server.ListenAndServe()
}
