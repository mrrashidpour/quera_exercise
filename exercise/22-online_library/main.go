package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Book struct {
	Title    string
	Author   string
	Borrowed bool
}

type Server struct {
	port  string
	books map[string]*Book
}

func NewServer(port string) *Server {
	return &Server{
		port:  port,
		books: make(map[string]*Book),
	}
}

func (s *Server) Start() {
	http.HandleFunc("/book", s.handleBook)
	http.ListenAndServe(":"+s.port, nil)
}

func key(title, author string) string {
	return strings.ToLower(title) + "::" + strings.ToLower(author)
}

func (s *Server) handleBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		s.handlePost(w, r)
	case http.MethodGet:
		s.handleGet(w, r)
	case http.MethodPut:
		s.handlePut(w, r)
	case http.MethodDelete:
		s.handleDelete(w, r)
	default:
		http.Error(w, `{"Result": "", "Error": "method not allowed"}`, http.StatusBadRequest)
	}
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := strings.ToLower(strings.TrimSpace(r.FormValue("title")))
	author := strings.ToLower(strings.TrimSpace(r.FormValue("author")))

	if title == "" || author == "" {
		http.Error(w, `{"Result": "", "Error": "title or author cannot be empty"}`, http.StatusBadRequest)
		return
	}

	k := key(title, author)
	if _, exists := s.books[k]; exists {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"Result": "this book is already in the library", "Error": ""}`)
		return
	}

	s.books[k] = &Book{Title: title, Author: author, Borrowed: false}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"Result": "added book %s by %s", "Error": ""}`, title, author)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := strings.ToLower(strings.TrimSpace(q.Get("title")))
	author := strings.ToLower(strings.TrimSpace(q.Get("author")))

	if title == "" || author == "" {
		http.Error(w, `{"Result": "", "Error": "title or author cannot be empty"}`, http.StatusBadRequest)
		return
	}

	k := key(title, author)
	book, exists := s.books[k]
	if !exists {
		http.Error(w, `{"Result": "", "Error": "this book does not exist"}`, http.StatusBadRequest)
		return
	}
	if book.Borrowed {
		http.Error(w, `{"Result": "", "Error": "this book is borrowed"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := strings.ToLower(strings.TrimSpace(q.Get("title")))
	author := strings.ToLower(strings.TrimSpace(q.Get("author")))

	if title == "" || author == "" {
		http.Error(w, `{"Result": "", "Error": "title or author cannot be empty"}`, http.StatusBadRequest)
		return
	}

	var body struct {
		Borrow *bool `json:"borrow"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Borrow == nil {
		http.Error(w, `{"Result": "", "Error": "borrow value cannot be empty"}`, http.StatusBadRequest)
		return
	}

	k := key(title, author)
	book, exists := s.books[k]
	if !exists {
		http.Error(w, `{"Result": "", "Error": "this book does not exist"}`, http.StatusBadRequest)
		return
	}

	if *body.Borrow {
		if book.Borrowed {
			http.Error(w, `{"Result": "", "Error": "this book is already borrowed"}`, http.StatusBadRequest)
			return
		}
		book.Borrowed = true
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"Result": "you have borrowed this book successfully", "Error": ""}`)
		return
	}

	if !book.Borrowed {
		http.Error(w, `{"Result": "", "Error": "this book is already in the library"}`, http.StatusBadRequest)
		return
	}
	book.Borrowed = false
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"Result": "thank you for returning this book", "Error": ""}`)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := strings.ToLower(strings.TrimSpace(q.Get("title")))
	author := strings.ToLower(strings.TrimSpace(q.Get("author")))

	if title == "" || author == "" {
		http.Error(w, `{"Result": "", "Error": "title or author cannot be empty"}`, http.StatusBadRequest)
		return
	}

	k := key(title, author)
	if _, exists := s.books[k]; !exists {
		http.Error(w, `{"Result": "", "Error": "this book does not exist"}`, http.StatusBadRequest)
		return
	}

	delete(s.books, k)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"Result": "successfully deleted", "Error": ""}`)
}
