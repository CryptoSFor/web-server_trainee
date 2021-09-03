package handlers

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"server/storage"
)

func HandleRequests() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/books", returnBooksHandler)
	router.HandleFunc("/book/{id}", returnBookHandler)
	router.HandleFunc("/edit/{id}", editPageHandler).Methods("GET")
	router.HandleFunc("/edit/{id}", editBookHandler).Methods("POST")
	router.HandleFunc("/delete/{id}", deleteBookHandler)
	return router
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

func returnBooksHandler(w http.ResponseWriter, r *http.Request) {
	db := storage.OpenDb()
	books := storage.ReturnAllBooks(db)
	tmpl, _ := template.ParseFiles("templates/books.html")
	err := tmpl.Execute(w, books)
	if err != nil {
		log.Println(err)
	}
}

func returnBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	db := storage.OpenDb()
	book := storage.ReturnSingleBook(db, key)
	if book.Id == 0{
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	}else {
		tmpl, _ := template.ParseFiles("templates/book.html")
		err := tmpl.Execute(w, book)
		if err != nil {
			log.Println(err)
		}
	}
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	db := storage.OpenDb()
	book := storage.ReturnSingleBook(db, key)
	if book.Id == 0 {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit.html")
		err := tmpl.Execute(w, book)
		if err != nil {
			log.Println(err)
		}
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
