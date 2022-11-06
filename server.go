package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/frontmatter"
	_ "github.com/glebarez/go-sqlite"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/mux"
)

type server struct {
	router    *mux.Router
	blogposts []Blogpost
	tmpl      *template.Template
}

type Blogpost struct {
	Url        string
	Title      string
	Content    template.HTML
	Date       time.Time
	DateString string
	Archive    bool
}

func parseBlog(url string, file *os.File) (Blogpost, error) {
	var matter struct {
		Title   string
		Date    string
		Archive bool
	}

	content, err := frontmatter.Parse(file, &matter)
	if err != nil {
		return Blogpost{}, err
	}

	extensions := parser.CommonExtensions | parser.FencedCode
	parser := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(content, parser, nil)

	parsedDate, err := time.Parse("02/01/2006", matter.Date)
	if err != nil {
		return Blogpost{}, err
	}

	return Blogpost{
		Url:        url,
		Title:      matter.Title,
		Content:    template.HTML(html),
		Date:       parsedDate,
		DateString: matter.Date,
		Archive:    matter.Archive,
	}, nil
}

func newServer() (*server, error) {
	r := mux.NewRouter()

	blogposts := make([]Blogpost, 0)

	files, err := ioutil.ReadDir("blog")
	if err != nil {
		return nil, err
	}

	for _, info := range files {
		filename := info.Name()
		var url string
		if len(filename) > 3 {
			url = filename[:len(filename)-3]
		}
		file, err := os.Open(filepath.Join("blog", filename))
		if err != nil {
			return nil, err
		}

		blogpost, err := parseBlog(url, file)
		if err != nil {
			return nil, err
		}

		if blogpost.Archive == false {
			blogposts = append(blogposts, blogpost)
		}
	}

	tmpl, err := template.ParseFiles("templates/home.gohtml", "templates/blog.gohtml")
	if err != nil {
		return nil, err
	}

	s := &server{router: r, blogposts: blogposts, tmpl: tmpl}
	s.routes()
	return s, nil
}
