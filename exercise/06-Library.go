package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Book struct {
	ISBN  int
	Title string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	books := make(map[int]string)

	for i := 0; i < n; i++ {
		scanner.Scan()
		line := scanner.Text()
		parts := strings.Split(line, " ")

		command := parts[0]

		switch command {
		case "ADD":
			isbn, _ := strconv.Atoi(parts[1])
			title := strings.Join(parts[2:], " ")
			books[isbn] = title

		case "REMOVE":
			isbn, _ := strconv.Atoi(parts[1])
			delete(books, isbn)
		}
	}

	bookList := make([]Book, 0, len(books))
	for isbn, title := range books {
		bookList = append(bookList, Book{ISBN: isbn, Title: title})
	}

	sort.Slice(bookList, func(i, j int) bool {
		if bookList[i].Title == bookList[j].Title {
			return bookList[i].ISBN < bookList[j].ISBN
		}
		return bookList[i].Title < bookList[j].Title
	})

	for _, book := range bookList {
		fmt.Println(book.ISBN)
	}
}
