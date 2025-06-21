namespace go book
include "model.thrift"

struct AddBookRequest{
    1: required string ISBN,
    2: required string location,
    3: required string status,
    4: required i64 purchase_date,
    5: required double purchase_price,
}
struct AddBookResponse{
    1: model.BaseResp base,
    2: required i64 book_id,
}

struct UpdateBookRequest{
    1:required i64 book_id,
    2: optional string location,
    3: optional string status,
    4: optional i64 purchase_date,
    5: optional double purchase_price,
}
struct UpdateBookResponse{
    1: model.BaseResp base,
    2: required model.Book data,
}

struct DeleteBookRequest{
    1: required i64 book_id,
}
struct DeleteBookResponse{
    1: model.BaseResp base,
}

struct GetBookRequest{
    1: optional i64 book_id,
    2: optional string ISBN,
    3: required i64 page_size,
    4: required i64 page_num,
}
struct GetBookResponse{
    1: model.BaseResp base,
    2: required list<model.Book> data,
    3: required i64 total_count,
}

service BookService {
    AddBookResponse addBook(1: AddBookRequest req)(api.post="/book/add"),
    UpdateBookResponse updateBook(1: UpdateBookRequest req)(api.put="/book/update"),
    DeleteBookResponse deleteBook(1: DeleteBookRequest req)(api.delete="/book/delete"),
    GetBookResponse getBook(1: GetBookRequest req)(api.get="/book/search"),
}
