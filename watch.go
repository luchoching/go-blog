package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join("templates", "Post.html")))
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
