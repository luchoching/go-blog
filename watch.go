package main

import (
	//"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

type Post struct {
	Title string
	Body  template.HTML
}

type Index struct {
	PostList []ListItem
}

// index list item 增加modified date, 用內容第1行當作title
type ListItem struct {
	Name    string
	ModTime string
}

func loadPost(title string) *Post {
	fileanme := title + ".md"
	source, _ := ioutil.ReadFile(path.Join("posts", fileanme))
	body := template.HTML(blackfriday.MarkdownCommon(source))
	return &Post{title, body}
}

func loadListItem(f os.FileInfo) ListItem {
	name := f.Name()
	t := f.ModTime()
	return ListItem{
		Name:    name[:len(name)-3],
		ModTime: t.Format("Mon Jan _2 15:04:05 2006"),
	}
}

func loadPostList() *Index {
	files, _ := ioutil.ReadDir("posts")
	l := []ListItem{}
	for _, f := range files {
		l = append(l, loadListItem(f))
	}
	return &Index{l}
}

var templates = map[string]*template.Template{
	"post": template.Must(template.ParseFiles(
		path.Join("templates", "Base.html"),
		path.Join("templates", "Sidebar.html"),
		path.Join("templates", "Post.html"),
	)),
	"index": template.Must(template.ParseFiles(
		path.Join("templates", "Base.html"),
		path.Join("templates", "Sidebar.html"),
		path.Join("templates", "Index.html"),
	)),
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := loadPost(r.URL.Path[len("/posts/"):])
	templates["post"].ExecuteTemplate(w, "base", p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates["index"].ExecuteTemplate(w, "base", loadPostList())
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/posts/", handler)
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
