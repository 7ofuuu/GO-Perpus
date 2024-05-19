package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Book struct represents the details of a book
type Book struct {
	Title     string
	Author    string
	ISBN      string
	Copies    int
	Borrowers []Borrower // To keep track of who borrowed the book and the due date
}

// Borrower struct represents a borrower with name and due date
type Borrower struct {
	Name    string
	DueDate time.Time
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

// BorrowBook borrows a book from the library
func (lib *Library) BorrowBook(isbn string, borrower string) {
	for i, book := range lib.Books {
		if book.ISBN == isbn {
			if book.Copies > 0 {
				if lib.CanBorrowMoreBooks(borrower) {
					if !lib.HasBorrowedSameBook(isbn, borrower) {
						dueDate := time.Now().AddDate(0, 0, 7)
						lib.Books[i].Copies--
						lib.Books[i].Borrowers = append(lib.Books[i].Borrowers, Borrower{Name: borrower, DueDate: dueDate})
						fmt.Printf("Book borrowed successfully by %s. Due date: %s\n", borrower, dueDate.Format("02-Jan-2006"))
					} else {
						fmt.Println("You have already borrowed this book.")
					}
				} else {
					fmt.Println("You cannot borrow more than 3 books.")
				}
			} else {
				fmt.Println("No copies available for borrowing.")
			}
			return
		}
	}
	fmt.Println("Book not found.")
}

// ReturnBook returns a book to the library
func (lib *Library) ReturnBook(isbn string, borrower string) {
	for i, book := range lib.Books {
		if book.ISBN == isbn {
			for j, b := range book.Borrowers {
				if b.Name == borrower {
					lib.Books[i].Copies++
					lib.Books[i].Borrowers = append(lib.Books[i].Borrowers[:j], lib.Books[i].Borrowers[j+1:]...)
					fmt.Printf("Book returned successfully by %s\n", borrower)
					return
				}
			}
		}
	}
	fmt.Println("Book was not borrowed by this user.")
}

// CanBorrowMoreBooks checks if the user has borrowed less than 3 books
func (lib *Library) CanBorrowMoreBooks(borrower string) bool {
	count := 0
	for _, book := range lib.Books {
		for _, b := range book.Borrowers {
			if b.Name == borrower {
				count++
				if count >= 3 {
					return false
				}
			}
		}
	}
	return true
}

// HasBorrowedSameBook checks if the user has already borrowed the same book
func (lib *Library) HasBorrowedSameBook(isbn string, borrower string) bool {
	for _, book := range lib.Books {
		if book.ISBN == isbn {
			for _, b := range book.Borrowers {
				if b.Name == borrower {
					return true
				}
			}
		}
	}
	return false
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

	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Edit Book")
		fmt.Println("4. Search Books")
		fmt.Println("5. Display All Books")
		fmt.Println("6. Borrow Book")
		fmt.Println("7. Return Book")
		fmt.Println("8. Exit")

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
			// Save library data to file after adding a new book
			if err := library.SaveLibrary("library.txt"); err != nil {
				fmt.Println("Error saving library data:", err)
			}
		case "2":
			isbn := getInput("Enter ISBN of the book to remove:")
			library.RemoveBook(isbn)
			fmt.Println("Book removed successfully.")
			// Save library data to file after removing a book
			if err := library.SaveLibrary("library.txt"); err != nil {
				fmt.Println("Error saving library data:", err)
			}
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
			// Save library data to file after editing a book
			if err := library.SaveLibrary("library.txt"); err != nil {
				fmt.Println("Error saving library data:", err)
			}
		case "4":
			keyword := getInput("Enter title or author to search:")
			searchResults := library.SearchBooks(keyword)
			fmt.Println("Search results:")
			displayBooks(searchResults)
		case "5":
			fmt.Println("All Books:")
			displayBooks(library.Books)
		case "6":
			fmt.Println("Books available for borrowing:")
			displayBooks(library.Books)
			isbn := getInput("Enter ISBN of the book to borrow:")
			borrower := getInput("Enter your name:")
			library.BorrowBook(isbn, borrower)
			// Save library data to file after borrowing a book
			if err := library.SaveLibrary("library.txt"); err != nil {
				fmt.Println("Error saving library data:", err)
			}
		case "7":
			isbn := getInput("Enter ISBN of the book to return:")
			borrower := getInput("Enter your name:")
			library.ReturnBook(isbn, borrower)
			// Save library data to file after returning a book
			if err := library.SaveLibrary("library.txt"); err != nil {
				fmt.Println("Error saving library data:", err)
			}
		case "8":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 8.")
		}
	}
}

// displayBooks displays a list of books in a table format
func displayBooks(books []Book) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Title", "Author", "ISBN", "Copies", "Borrowers"})
	for _, book := range books {
		var borrowerDetails []string
		for _, borrower := range book.Borrowers {
			borrowerDetails = append(borrowerDetails, fmt.Sprintf("%s (Due: %s)", borrower.Name, borrower.DueDate.Format("02-Jan-2006")))
		}
		borrowers := strings.Join(borrowerDetails, ", ")
		table.Append([]string{book.Title, book.Author, book.ISBN, fmt.Sprintf("%d", book.Copies), borrowers})
	}
	table.Render()
}
