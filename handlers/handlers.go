package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"server/storage"
)

func HandleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/books", booksHandler)
	router.HandleFunc("/book/{id}", bookHandler)
	router.HandleFunc("/edit/{id}", editPageHandler).Methods("GET")
	router.HandleFunc("/edit/{id}", editBookHandler).Methods("POST")
	router.HandleFunc("/delete/{id}", deleteBookHandler)
	log.Fatal(http.ListenAndServe(":8081", router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		name := r.FormValue("name")
		author := r.FormValue("author")
		genre := r.FormValue("genre")
		db := storage.OpenDb()
		storage.AddBook(db, name, author, genre)
		http.Redirect(w, r, "/books", 301)
	}else {
		http.ServeFile(w,r, "templates/index.html")
	}
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	db := storage.OpenDb()
	books := storage.ReturnAllBooks(db)
	tmpl, _ := template.ParseFiles("templates/books.html")
	err := tmpl.Execute(w, books)
	if err != nil {
		log.Println(err)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	db := storage.OpenDb()
	book := storage.ReturnSingleBook(db, key)
	if book.Id == 0{
		json.NewEncoder(w).Encode("NULL")
	}else {
		json.NewEncoder(w).Encode(book)
	}
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	db := storage.OpenDb()
	book := storage.ReturnSingleBook(db, key)
	tmpl, _ := template.ParseFiles("templates/edit.html")
	err := tmpl.Execute(w, book)
	if err != nil {
		log.Println(err)
	}
}

func editBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	name := r.FormValue("name")
	author := r.FormValue("author")
	genre := r.FormValue("genre")
	db := storage.OpenDb()
	storage.UpdateBook(db, key, name, author, genre)
	http.Redirect(w, r, "/books", 301)
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	db := storage.OpenDb()
	storage.DeleteBook(db, key)
	http.Redirect(w, r, "/books", 301)
}