package main

import (
	"html/template"
	"log"
	"net/http"
)

type Film struct {
	Title    string
	Director string
}

func main() {

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		films := map[string][]Film{
			"Films": {
				{Title: "The GodFather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}

		tmpl.Execute(w, films)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		/* htmlStr:=fmt.Sprintf("<li class='list-group-item'> %s - %s </li>",title,director)
		tmpl,_:=template.New("t").Parse(htmlStr)
		tmpl.Execute(w,nil)	 */ //!köylü işi
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})

	}
	h3 := func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w,"film-list-element",id)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-film", h2)
	http.HandleFunc("/delete", h3)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
