package db

import (
	"time"

	"github.com/2451965602/LMS/pkg/constants"
)

type User struct {
	ID           int64     `json:"id"            gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name"          gorm:"type:varchar(50);not null;unique"`
	Password     string    `json:"password"      gorm:"type:varchar(255);not null"`
	Permission   string    `json:"permissions"   gorm:"type:enum('admin','librarian','member');default:'member';not null"`
	Phone        *string   `json:"phone"         gorm:"type:varchar(20)"`
	RegisterDate time.Time `json:"register_date" gorm:"column:pegister_date;type:timestamp;default:CURRENT_TIMESTAMP;not null"`
	Status       string    `json:"status"        gorm:"type:enum('active','suspended','inactive');default:'active';not null"`
}

func (User) TableName() string {
	return constants.UserTableName
}

type BookType struct {
	ISBN            string `json:"isbn"             gorm:"type:varchar(20);primaryKey"`
	Title           string `json:"title"            gorm:"type:varchar(100);not null"`
	Author          string `json:"author"           gorm:"type:varchar(50);not null"`
	Category        string `json:"category"         gorm:"type:varchar(50);not null"`
	Publisher       string `json:"publisher"        gorm:"type:varchar(50);not null"`
	PublishYear     int64  `json:"publish_year"     gorm:"type:int;not null"`
	Description     string `json:"description"      gorm:"type:text"`
	TotalCopies     int64  `json:"total_copies"     gorm:"type:int;default:0;not null"`
	AvailableCopies int64  `json:"available_copies" gorm:"type:int;default:0;not null"`
}

func (BookType) TableName() string {
	return constants.BookTypeTableName
}

type Book struct {
	ID            int64      `json:"id"             gorm:"primaryKey;autoIncrement"`
	ISBN          string     `json:"isbn"           gorm:"type:varchar(20);not null"`
	Location      string     `json:"location"       gorm:"type:varchar(50);not null"`
	Status        string     `json:"status"         gorm:"type:enum('available','checked_out','lost','damaged');default:'available'"`
	PurchaseDate  time.Time  `json:"purchase_date"  gorm:"type:timestamp;not null"`
	PurchasePrice float64    `json:"purchase_price" gorm:"type:decimal(10,2);not null"`
	LastCheckout  *time.Time `json:"last_checkout"  gorm:"type:timestamp"`
}

func (Book) TableName() string {
	return constants.BookTableName
}

type BorrowRecord struct {
	ID           int64      `json:"id"             gorm:"primaryKey;autoIncrement"`
	UserID       int64      `json:"user_id"        gorm:"not null"`
	BookID       int64      `json:"book_id"        gorm:"not null"`
	CheckoutDate time.Time  `json:"checkout_date"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	RenewalCount int64      `json:"renewal_count"  gorm:"type:int;default:0"`
	DueDate      time.Time  `json:"due_date"       gorm:"type:timestamp;not null"`
	ReturnDate   *time.Time `json:"return_date"    gorm:"type:timestamp"`
	Status       string     `json:"status"         gorm:"type:enum('checked_out','returned','overdue','lost');default:'checked_out'"`
	LateFee      float64    `json:"late_fee"       gorm:"type:decimal(10,2);default:0.00"`
}

func (BorrowRecord) TableName() string {
	return constants.BorrowRecordTableName
}

type Reservation struct {
	ID          int64     `json:"id"             gorm:"primaryKey;autoIncrement"`
	UserID      int64     `json:"user_id"        gorm:"not null"`
	BookID      int64     `json:"book_id"        gorm:"not null"`
	ISBN        string    `json:"isbn"           gorm:"type:varchar(20);not null"`
	ReserveDate time.Time `json:"reserve_date"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	ExpiryDate  time.Time `json:"expiry_date"    gorm:"type:timestamp;not null"`
	Status      string    `json:"status"         gorm:"type:enum('pending','fulfilled','cancelled','expired');default:'pending'"`
}

func (Reservation) TableName() string {
	return constants.ReservationTableName
}
