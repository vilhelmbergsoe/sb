package main

import (
	"net/http"
	"sort"
	"io/fs"

	"github.com/gorilla/mux"
)

func (s *server) routes() error {
	fsys := fs.FS(content)
	static, err := fs.Sub(fsys, "static")
	if err != nil {
		return err
	}
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	s.router.HandleFunc("/", s.HandleHome()).Methods("GET")
	s.router.HandleFunc("/blog/{url}", s.HandleBlogpost()).Methods("GET")

	return nil
}

func (s *server) HandleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sort.Slice(s.blogposts, func(i, j int) bool {
			return s.blogposts[i].Date.After(s.blogposts[j].Date)
		})

		s.tmpl.ExecuteTemplate(w, "home.gohtml", s.blogposts)
	}
}

func (s *server) HandleBlogpost() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blogpostUrl := vars["url"]

		var blogpost Blogpost

		for _, v := range s.blogposts {
			if blogpostUrl == v.Url {
				blogpost = v
			}
		}
		if blogpost.Title == "" {
			http.Error(w, "404 Blog post not found", http.StatusNotFound)
			return
		}

		s.tmpl.ExecuteTemplate(w, "blog.gohtml", blogpost)
	}
}
