package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (s *server) routes() {
	s.router.HandleFunc("/", s.HandleHome()).Methods("GET")
	s.router.HandleFunc("/admin", s.Auth(s.HandleAdmin())).Methods("GET")
	s.router.HandleFunc("/blogposts", s.HandleBlogposts()).Methods("GET")
	s.router.HandleFunc("/blogposts", s.Auth(s.HandleCreateBlogpost())).Methods("POST")
	s.router.HandleFunc("/blogposts/delete/{id}", s.Auth(s.HandleDeleteBlogpost())).Methods("POST")
	s.router.HandleFunc("/blogposts/update/{id}", s.Auth(s.HandleUpdateBlogpost())).Methods("POST")
}

func (s *server) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if ok {
			usernameSha256Hash := sha256.New()
			passwordSha256Hash := sha256.New()
			io.Copy(usernameSha256Hash, strings.NewReader(username))
			io.Copy(passwordSha256Hash, strings.NewReader(password))

			userHash := fmt.Sprintf("%x", usernameSha256Hash.Sum(nil))
			passHash := fmt.Sprintf("%x", passwordSha256Hash.Sum(nil))

			row, err := s.db.Query("SELECT * FROM users;")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer row.Close()

			users := make([]User, 0)
			for row.Next() {
				var user User
				row.Scan(&user.ID, &user.Username, &user.Password)
				users = append(users, user)
			}

			if err = row.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, user := range users {
				if user.Username == userHash && user.Password == passHash {
					h(w, r)
					return
				}
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (s *server) HandleHome() http.HandlerFunc {
	type Data struct {
		Blogposts []Blogpost
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:8080/blogposts")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		buf := bufio.NewReader(resp.Body)

		body, err := ioutil.ReadAll(buf)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		var data Data

		json.Unmarshal(body, &data.Blogposts)

		s.hometmpl.Execute(w, data)
	}
}

func (s *server) HandleAdmin() http.HandlerFunc {

	type Data struct {
		Blogposts []Blogpost
	}

	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:8080/blogposts")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		buf := bufio.NewReader(resp.Body)

		body, err := ioutil.ReadAll(buf)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		var data Data

		json.Unmarshal(body, &data.Blogposts)

		s.admintmpl.Execute(w, data)
	}
}

func (s *server) HandleBlogposts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blogposts []Blogpost

		row, err := s.db.Query("SELECT * FROM blogposts ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer row.Close()
		for row.Next() {
			var blogpost Blogpost
			if err = row.Scan(&blogpost.ID, &blogpost.Title, &blogpost.Content); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			blogposts = append(blogposts, blogpost)
		}

		if err = row.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blogposts)
	}
}

func (s *server) HandleDeleteBlogpost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deleteBlogpost := `DELETE FROM blogposts WHERE id = ?`

		stmnt, err := s.db.Prepare(deleteBlogpost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = stmnt.Exec(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) HandleUpdateBlogpost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var blogpost Blogpost

		err = json.NewDecoder(r.Body).Decode(&blogpost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updateBlogpost := `UPDATE blogposts SET title = ?, content = ? WHERE id = ?`

		stmnt, err := s.db.Prepare(updateBlogpost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = stmnt.Exec(blogpost.Title, blogpost.Content, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) HandleCreateBlogpost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blogpost Blogpost

		err := json.NewDecoder(r.Body).Decode(&blogpost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		insertBlogpost := `INSERT INTO blogposts VALUES (NULL, ?, ?)`

		stmnt, err := s.db.Prepare(insertBlogpost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = stmnt.Exec(blogpost.Title, blogpost.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
