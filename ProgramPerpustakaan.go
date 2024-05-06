package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Book struct represents the details of a book
type Book struct {
	Title  string
	Author string
	ISBN   string
	Copies int
	// Add other book details as needed
}

// Library struct represents the library with its collection of books
type Library struct {
	Books []Book
}

// AddBook adds a new book to the library
func (lib *Library) AddBook(book Book) {
	lib.Books = append(lib.Books, book)
}

// SearchBooks searches for books by title or author
func (lib *Library) SearchBooks(keyword string) []Book {
	var results []Book
	for _, book := range lib.Books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(keyword)) || strings.Contains(strings.ToLower(book.Author), strings.ToLower(keyword)) {
			results = append(results, book)
		}
	}
	return results
}

// EditBook edits the details of a book in the library
func (lib *Library) EditBook(isbn string, newDetails Book) {
	for i, book := range lib.Books {
		if book.ISBN == isbn {
			lib.Books[i] = newDetails
			break
		}
	}
}

// RemoveBook removes a book from the library
func (lib *Library) RemoveBook(isbn string) {
	for i, book := range lib.Books {
		if book.ISBN == isbn {
			lib.Books = append(lib.Books[:i], lib.Books[i+1:]...)
			break
		}
	}
}

// SortBooks sorts the library's books by title
func (lib *Library) SortBooks() {
	sort.Slice(lib.Books, func(i, j int) bool {
		return lib.Books[i].Title < lib.Books[j].Title
	})
}

// SaveLibrary saves the library data to a file
func (lib *Library) SaveLibrary(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, book := range lib.Books {
		_, err := fmt.Fprintf(file, "%s,%s,%s,%d\n", book.Title, book.Author, book.ISBN, book.Copies)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadLibrary loads the library data from a file
func (lib *Library) LoadLibrary(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var books []Book
	var title, author, isbn string
	var copies int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) != 4 {
			continue
		}
		title, author, isbn = fields[0], fields[1], fields[2]
		fmt.Sscanf(fields[3], "%d", &copies)
		books = append(books, Book{Title: title, Author: author, ISBN: isbn, Copies: copies})
	}

	lib.Books = books
	return nil
}

// getInput is a helper function to get user input
func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	// Initialize library
	library := Library{}

	// Load library data from file
	if err := library.LoadLibrary("library.txt"); err != nil {
		fmt.Println("Error loading library data:", err)
	}

	// Add some sample books if the library is empty
	if len(library.Books) == 0 {
		for i := 0; i < 3; i++ {
			fmt.Println("Enter details for book", i+1)
			title := getInput("Title:")
			author := getInput("Author:")
			isbn := getInput("ISBN:")
			copies := getInput("Number of copies:")
			var copiesInt int
			fmt.Sscanf(copies, "%d", &copiesInt)
			library.AddBook(Book{Title: title, Author: author, ISBN: isbn, Copies: copiesInt})
		}
	}

	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Edit Book")
		fmt.Println("4. Search Books")
		fmt.Println("5. Display All Books")
		fmt.Println("6. Exit")

		choice := getInput("Enter your choice:")

		switch choice {
		case "1":
			fmt.Println("Adding a new book:")
			title := getInput("Title:")
			author := getInput("Author:")
			isbn := getInput("ISBN:")
			copies := getInput("Number of copies:")
			var copiesInt int
			fmt.Sscanf(copies, "%d", &copiesInt)
			library.AddBook(Book{Title: title, Author: author, ISBN: isbn, Copies: copiesInt})
			fmt.Println("Book added successfully.")
		case "2":
			isbn := getInput("Enter ISBN of the book to remove:")
			library.RemoveBook(isbn)
			fmt.Println("Book removed successfully.")
		case "3":
			isbn := getInput("Enter ISBN of the book to edit:")
			newTitle := getInput("New Title:")
			newAuthor := getInput("New Author:")
			newISBN := getInput("New ISBN:")
			newCopies := getInput("New Number of copies:")
			var newCopiesInt int
			fmt.Sscanf(newCopies, "%d", &newCopiesInt)
			library.EditBook(isbn, Book{Title: newTitle, Author: newAuthor, ISBN: newISBN, Copies: newCopiesInt})
			fmt.Println("Book edited successfully.")
		case "4":
			keyword := getInput("Enter title or author to search:")
			searchResults := library.SearchBooks(keyword)
			fmt.Println("Search results:")
			displayBooks(searchResults)
		case "5":
			fmt.Println("All Books:")
			displayBooks(library.Books)
		case "6":
			fmt.Println("Exiting...")
			// Save library data to file before exiting
			if err := library.SaveLibrary("library.txt"); err != nil {
				fmt.Println("Error saving library data:", err)
			}
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 6.")
		}
	}
}

// displayBooks displays a list of books in a table format
func displayBooks(books []Book) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "Author", "ISBN", "Copies"})
	for _, book := range books {
		table.Append([]string{book.Title, book.Author, book.ISBN, fmt.Sprintf("%d", book.Copies)})
	}
	table.Render()
}
