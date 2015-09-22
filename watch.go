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

type Index struct {
	PostList []string
}

func loadPost(title string) *Post {
	fileanme := title + ".md"
	source, _ := ioutil.ReadFile(path.Join("posts", fileanme))
	body := template.HTML(blackfriday.MarkdownCommon(source))
	return &Post{title, body}
}

func loadPostList() *Index {
	files, _ := ioutil.ReadDir("posts")
	l := []string{}
	for _, f := range files {
		name := f.Name()
		l = append(l, name[:len(name)-3])
	}
	return &Index{l}
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := loadPost(r.URL.Path[len("/posts/"):])
	t := template.Must(template.ParseFiles(path.Join("templates", "Post.html")))
	t.Execute(w, p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(path.Join("templates", "Index.html")))
	t.Execute(w, loadPostList())
}

func main() {
	http.HandleFunc("/posts/", handler)
	http.HandleFunc("/", indexHandler)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
