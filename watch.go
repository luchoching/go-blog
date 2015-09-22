package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Post struct {
	Title string
	Body  string
}

func loadPost(title string) *Post {
	fileanme := title + ".md"
	body, _ := ioutil.ReadFile(path.Join("posts", fileanme))
	fmt.Println(fileanme)
	fmt.Println(body)
	return &Post{title, string(body)}
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := loadPost(r.URL.Path[1:])
	t := template.Must(template.ParseFiles(path.Join("templates", "Post.html")))
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
