package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Id      int
	Content string
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		return db, fmt.Errorf("Error opening db: %v", err)
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS entries " +
			"(id INTEGER PRIMARY KEY, entry VARCHAR(256));")
	if err != nil {
		return db, fmt.Errorf("Error creating table: %v", err)
	}

	log.Printf("Database connected and setup")
	return db, nil
}

func getTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT id, entry FROM entries;")
	if err != nil {
		return nil, fmt.Errorf("Error getting todos: %v", err)
	}
	todos := make([]Todo, 0)
	for rows.Next() {
		var id int
		var content string
		err = rows.Scan(&id, &content)
		if err != nil {
			return nil, fmt.Errorf("Couldn't get posts from db: %v", err)
		}
		todos = append(todos, Todo{Id: id, Content: content})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Row error: %v", err)
	}
	return todos, nil
}

func addTodo(db *sql.DB, tt Todo) error {
	_, err := db.Exec(
		"INSERT INTO entries (entry)"+
			"VALUES (?);",
		tt.Content,
	)
	if err != nil {
		return fmt.Errorf("Error inserting post: %v", err)
	}
	return nil
}

func deleteTodo(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM entries WHERE id = ?;", id)
	if err != nil {
		return fmt.Errorf("Error inserting post: %v", err)
	}
	return nil
}

type THandler struct {
	db *sql.DB
	ts *template.Template
}

func (hh THandler) ServeHTTP(ww http.ResponseWriter, rr *http.Request) {
	var info struct {
		Todos []Todo
	}
	var err error
	info.Todos, err = getTodos(hh.db)
	if err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := hh.ts.ExecuteTemplate(
		ww,
		"main.html",
		info); err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
	}
}

type PHandler struct {
	db *sql.DB
}

func (hh PHandler) ServeHTTP(ww http.ResponseWriter, rr *http.Request) {
	tt := Todo{Content: rr.FormValue("content")}
	log.Printf("New todo: %v", tt.Content)
	if err := addTodo(hh.db, tt); err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(ww, rr, "/", http.StatusFound)
}

type DHandler struct {
	db *sql.DB
}

func (hh DHandler) ServeHTTP(ww http.ResponseWriter, rr *http.Request) {
	id, err := strconv.Atoi(rr.FormValue("id"))
	if err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := deleteTodo(hh.db, id); err != nil {
		http.Error(ww, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(ww, rr, "/", http.StatusFound)
}

func main() {
	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	ts := template.Must(template.ParseFiles("main.html"))

	th := THandler{db: db, ts: ts}
	http.Handle("/", th)

	ph := PHandler{db: db}
	http.Handle("/add", ph)

	dh := DHandler{db: db}
	http.Handle("/delete", dh)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static/"))))

	http.ListenAndServe(":8080", nil)
}
