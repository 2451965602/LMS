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
}
struct ReturnResponse{
    1: model.BaseResp base,
    2: required model.BorrowRecord data,
}

struct RenewRequest{
    1: required i64 borrow_id,
}
struct RenewResponse{
    1: model.BaseResp base,
    2: required model.BorrowRecord data,
}

struct GetBorrowRecordRequest{
    1: optional i64 user_id,
    2: optional i64 book_id,
}
struct GetBorrowRecordResponse{
    1: model.BaseResp base,
    2: required list<model.BorrowRecord> data,
}

struct ReservationRequest{
    1: required i64 book_id,
    2: required string reserve_date,
}
struct ReservationResponse{
    1: model.BaseResp base,
    2: required i64 reservation_id,
}

struct CancelReservationRequest{
    1: required i64 reservation_id,
}
struct CancelReservationResponse{
    1: model.BaseResp base,
    2: required model.Reservation data,
}

struct GetReservationRequest{
    1: optional i64 user_id,
    2: optional i64 book_id,
}
struct GetReservationResponse{
    1: model.BaseResp base,
    2: required list<model.Reservation> data,
}

service BorrowService {
    BorrowResponse borrow(1: BorrowRequest req)(api.post="/book/borrow"),
    ReturnResponse returnBook(1: ReturnRequest req)(api.post="/book/return"),
    RenewResponse renew(1: RenewRequest req)(api.post="/book/renew"),
    GetBorrowRecordResponse getBorrowRecord(1: GetBorrowRecordRequest req)(api.get="/book/record"),
}

service ReservationService {
    ReservationResponse reserve(1: ReservationRequest req)(api.post="/reserve"),
    CancelReservationResponse cancelReservation(1: CancelReservationRequest req)(api.post="/reserve/cancel"),
    GetReservationResponse getReservation(1: GetReservationRequest req)(api.get="/reserve/record"),
}


