package main

import (
	"database/sql"
	"html/template"

	_ "github.com/glebarez/go-sqlite"
	"github.com/gorilla/mux"
)

type server struct {
	router       *mux.Router
	db           *sql.DB
	hometmpl     *template.Template
	admintmpl    *template.Template
	blogposttmpl *template.Template
}

type Blogpost struct {
	ID      int    `json:"id"` // gorm:"primary_key"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type User struct {
	ID       int    `json:"id"` // gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func newServer() (*server, error) {
	r := mux.NewRouter()

	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		return nil, err
	}

	createUserTable := `CREATE TABLE IF NOT EXISTS users(id integer primary key, u TEXT NOT NULL, p TEXT NOT NULL);`
	createBlogpostsTable := `CREATE TABLE IF NOT EXISTS blogposts(id integer primary key, title TEXT NOT NULL, content TEXT NOT NULL, date TEXT NOT NULL);`
	stmnt, err := db.Prepare(createUserTable)
	if err != nil {
		return nil, err
	}
	stmnt.Exec()

	stmnt, err = db.Prepare(createBlogpostsTable)
	if err != nil {
		return nil, err
	}
	stmnt.Exec()

	// err = db.AutoMigrate(&Blogpost{}, &User{})
	// if err != nil {
	// 	return nil, err
	// }

	hometmpl, err := template.ParseFiles("public/index.html")
	if err != nil {
		return nil, err
	}

	admintmpl, err := template.ParseFiles("public/admin.html")
	if err != nil {
		return nil, err
	}

	blogposttmpl, err := template.ParseFiles("public/blogpost.html")
	if err != nil {
		return nil, err
	}

	s := &server{router: r, db: db, hometmpl: hometmpl, admintmpl: admintmpl, blogposttmpl: blogposttmpl}
	s.routes()
	return s, nil
}
