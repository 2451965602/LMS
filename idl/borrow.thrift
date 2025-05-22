namespace go borrow
include "model.thrift"

struct BorrowRequest{
    1: required i64 book_id,
}
struct BorrowResponse{
    1: model.BaseResp base,
    2: required i64 borrow_id,
}

struct ReturnRequest{
    1: required i64 borrow_id,
    2: required i64 book_id,
    3: required string status,
    4: required double late_fee,
}
struct ReturnResponse{
    1: model.BaseResp base,
    2: required model.BorrowRecord data,
}

struct RenewRequest{
    1: required i64 borrow_id,
    2: required i64 add_time,
}
struct RenewResponse{
    1: model.BaseResp base,
    2: required model.BorrowRecord data,
}

struct GetBorrowRecordRequest{
    1: required i64 user_id,
    2: required i64 page_size,
    3: required i64 page_num,
    4: required i64 status,
}
struct GetBorrowRecordResponse{
    1: model.BaseResp base,
    2: required list<model.BorrowRecord> data,
    3: required i64 total,
}


service BorrowService {
    BorrowResponse borrow(1: BorrowRequest req)(api.post="/book/borrow"),
    ReturnResponse returnBook(1: ReturnRequest req)(api.post="/book/return"),
    RenewResponse renew(1: RenewRequest req)(api.post="/book/renew"),
    GetBorrowRecordResponse getBorrowRecord(1: GetBorrowRecordRequest req)(api.get="/book/record"),
}


