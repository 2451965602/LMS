package service

import "time"

type User struct {
	ID           int64
	Name         string
	Password     string
	Permissions  string
	Phone        *string
	RegisterDate time.Time
	Status       string
}

type BookType struct {
	ISBN            string
	Title           string
	Author          string
	Category        string
	Publisher       string
	PublishYear     *int
	Description     string
	TotalCopies     int
	AvailableCopies int
}

type Book struct {
	ID            int64
	ISBN          string
	Location      string
	Status        string
	PurchaseDate  *time.Time
	PurchasePrice *float64
	LastCheckout  *time.Time
	BookType      BookType
}

type BorrowRecord struct {
	ID           int64
	UserID       int64
	BookID       int64
	CheckoutDate time.Time
	DueDate      time.Time
	ReturnDate   *time.Time
	Status       string
	LateFee      float64
	User         User
	Book         Book
}

type Reservation struct {
	ID          int64
	UserID      int64
	ISBN        string
	ReserveDate time.Time
	ExpiryDate  time.Time
	Status      string
	User        User
	BookType    BookType
}
