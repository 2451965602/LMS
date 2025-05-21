namespace go model

struct BaseResp {
    1: required i64 code,
    2: required string msg,
}

struct User {
    1: required i64 id
    2: required string username
    3: required string password
    4: required string email
    5: required string phone
    6: required i32 status
    7: required string create_time
    8: required string update_time
}

struct BookType {
    1: required i64 id
    2: required string name
    3: required string description
    4: required string create_time
    5: required string update_time
}

struct Book {
    1: required i64 id
    2: required string title
    3: required string author
    4: required string isbn
    5: required i64 type_id
    6: required i32 status
    7: required string publish_date
    8: required string create_time
    9: required string update_time
}

struct BorrowRecord {
    1: required i64 id
    2: required i64 user_id
    3: required i64 book_id
    4: required string borrow_date
    5: required string due_date
    6: required string return_date
    7: required i32 status
    8: required string create_time
    9: required string update_time
}

struct Reservation {
    1: required i64 id
    2: required i64 user_id
    3: required i64 book_id
    4: required string reserve_date
    5: required string expire_date
    6: required i32 status
    7: required string create_time
    8: required string update_time
}