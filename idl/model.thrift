namespace go model

struct BaseResp {
    1: required i64 code,
    2: required string msg,
}

struct ErrorResp {
    1: BaseResp base,
}

struct User {
    1: required i64 id
    2: required string username
    3: required string password
    4: optional string phone
    5: required string status
    6: required string permissions
    7: required string register_date
}

struct BookType {
    1: required string ISBN
    2: required string title
    3: required string author
    4: required string category
    5: required string publisher
    6: required i64 publish_year
    7: required string description
    8: required i64 total_copies
    9: required i64 available_copies
}

struct Book {
    1: required i64 id
    2: required string isbn
    3: required string location
    4: required string status
    5: required string purchase_date
    6: required double purchase_price
    7: required string last_checkout

}

struct BorrowRecord {
    1: required i64 id
    2: required i64 user_id
    3: required i64 book_id
    4: required string title
    5: required string checkout_date
    6: required string due_date
    7: required string return_date
    8: required string status
    9: required i64 renewal_count
    10: required double late_fee
}


