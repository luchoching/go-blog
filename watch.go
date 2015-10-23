package main

import (
	"bufio"
	//"fmt"
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
	PostList  []ListItem
	DraftList []ListItem
}

type ListItem struct {
	URL     string
	ModTime string
	Title   string
}

func loadPost(srcDir string, title string) *Post {
	filename := title + ".md"
	source, _ := ioutil.ReadFile(path.Join(srcDir, filename))
	body := template.HTML(blackfriday.MarkdownCommon(source))
	return &Post{title, body}
}

func loadPostFirstLine(postType string, filename string) string {
	f, err := os.Open(path.Join(postType, filename))
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

func loadListItem(postType string, f os.FileInfo) ListItem {
	name := f.Name()
	t := f.ModTime()
	return ListItem{
		URL:     name[:len(name)-3],
		ModTime: t.Format("Mon Jan _2 15:04:05 2006"),
		Title:   loadPostFirstLine(postType, name),
	}
}

func loadPostList() *Index {
	postFiles, _ := ioutil.ReadDir("posts")
	draftFiles, _ := ioutil.ReadDir("drafts")
	p := []ListItem{}
	d := []ListItem{}
	for _, f := range postFiles {
		p = append(p, loadListItem("posts", f))
	}
	for _, f := range draftFiles {
		d = append(d, loadListItem("drafts", f))
	}

	return &Index{p, d}
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
	p := loadPost("posts", r.URL.Path[len("/post/"):])
	templates["post"].ExecuteTemplate(w, "base", p)
}

func draftHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPost("drafts", r.URL.Path[len("/draft/"):])
	templates["post"].ExecuteTemplate(w, "base", p)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	p := loadPost("pages", r.URL.Path[len("/page/"):])
	templates["post"].ExecuteTemplate(w, "base", p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates["index"].ExecuteTemplate(w, "base", loadPostList())
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/post/", handler)
	http.HandleFunc("/draft/", draftHandler)
	http.HandleFunc("/page/", pageHandler)
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
