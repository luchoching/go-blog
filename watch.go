package main

import (
	//"fmt"
	"bufio"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Post struct {
	Title string
	Body  template.HTML
}

type Index struct {
	PostList []ListItem
}

type ListItem struct {
	URL     string
	ModTime string
	Title   string
}

func loadPost(title string) *Post {
	fileanme := title + ".md"
	source, _ := ioutil.ReadFile(path.Join("posts", fileanme))
	body := template.HTML(blackfriday.MarkdownCommon(source))
	return &Post{title, body}
}

func loadPostFirstLine(filename string) string {
	f, err := os.Open(path.Join("posts", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var firstLine string
	for scanner.Scan() {
		firstLine = scanner.Text()
		break
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return strings.TrimPrefix(firstLine, "# ")
}

func loadListItem(f os.FileInfo) ListItem {
	name := f.Name()
	t := f.ModTime()
	return ListItem{
		URL:     name[:len(name)-3],
		ModTime: t.Format("Mon Jan _2 15:04:05 2006"),
		Title:   loadPostFirstLine(name),
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
		path.Join("templates", "Post.html"),
	)),
	"index": template.Must(template.ParseFiles(
		path.Join("templates", "Base.html"),
		path.Join("templates", "Index.html"),
	)),
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := loadPost(r.URL.Path[len("/post/"):])
	templates["post"].ExecuteTemplate(w, "base", p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates["index"].ExecuteTemplate(w, "base", loadPostList())
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/post/", handler)
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
