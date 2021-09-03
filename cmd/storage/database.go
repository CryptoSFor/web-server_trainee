package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "aboba"
	dbname   = "books"
)

type Book struct {
	Id     int
	Name   string
	Author string
	Genre  string
}

func OpenDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	return db
}

func ReturnAllBooks(db *sql.DB) []Book {
	rows, err := db.Query("select * from books order by id")
	if err != nil {
		log.Println(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
	var books []Book
	for rows.Next() {
		b := Book{}
		err := rows.Scan(&b.Id, &b.Name, &b.Author, &b.Genre)
		if err != nil {
			fmt.Println(err)
			continue
		}
		books = append(books, b)
	}
	return books
}

func ReturnSingleBook(db *sql.DB, key string) Book {
	row := db.QueryRow("select * from books where Id = $1", key)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
	b := Book{}
	err := row.Scan(&b.Id, &b.Name, &b.Author, &b.Genre)
	if err != nil {
		log.Println(err)
	}
	return b
}

func AddBook(db *sql.DB, name, author, genre string) {
	_, err := db.Exec("insert into books (name, author, genre) values ($1, $2, $3)",
		name, author, genre)
	if err != nil {
		log.Println(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
}

func UpdateBook(db *sql.DB, key, name, author, genre string) {
	_, err := db.Exec("update books set name = $1, author = $2, genre = $3 where id = $4",
		name, author, genre, key)
	if err != nil {
		log.Println(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
}

func DeleteBook(db *sql.DB, key string) {
	_, err := db.Exec("delete from books where id = $1", key)
	if err != nil {
		log.Println(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
}
