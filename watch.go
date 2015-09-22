package main

import (
	//"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Post struct {
	Title string
	Body  template.HTML
}

func loadPost(title string) *Post {
	fileanme := title + ".md"
	source, _ := ioutil.ReadFile(path.Join("posts", fileanme))
	body := template.HTML(blackfriday.MarkdownCommon(source))
	return &Post{title, body}
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := loadPost(r.URL.Path[len("/posts/"):])
	t := template.Must(template.ParseFiles(path.Join("templates", "Post.html")))
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/posts/", handler)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
