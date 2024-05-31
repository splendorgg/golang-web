package main

import (
	"html/template"
	"log"
	"net/http"
)

// Create Addres struct that contains addr as string
type Address struct {
	addr string
}

// NewAPIServer function takes our addres (:8080) from main.go -> server := NewAPIServer(":8080").
func NewAPIServer(addres string) *Address {
	return &Address{
		addr: addres,
	}
}

func (s *Address) Run() error {
	router := http.NewServeMux()                                                            //Creates new HTTP router.
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) { //Endpoint accepts GET requests and dynamically retrieves the {userID}
		userID := r.PathValue("userID")             // Gets {userID} from the URL
		tpl, _ := template.ParseFiles("index.html") // Loads the "index.html" and parses as an HTML template
		tpl.Execute(w, userID)                      // Executes the HTML template and passes the userID to use in HTML file
		w.Write([]byte("User id is: " + userID))
	})

	//Creates http.Server object, contains TCP Adress as "host:post" and how server will operate(handler)
	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}
	log.Printf("server has started at %s", s.addr)
	return server.ListenAndServe() //Startes the HTTP Server and listens request from specific address
}

func main() {
	server := NewAPIServer(":8080")
	server.Run()

}
