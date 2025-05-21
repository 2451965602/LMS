namespace go book
include "model.thrift"

struct AddBookRequest{
    1: required string ISBN,
    2: required string location,
    3: required string status,
    4: required string purchasedate,
    5: required string purchaseprice,
}
struct AddBookResponse{
    1: model.BaseResp base,
    2: required i64 book_id,
}

struct UpdateBookRequest{
    1:required i64 book_id,
    2: optional string location,
    3: optional string status,
    4: optional string purchasedate,
    5: optional string purchaseprice,
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
}
struct GetBookResponse{
    1: model.BaseResp base,
    2: required list<model.Book> data,
}

service BookService {
    AddBookResponse addBook(1: AddBookRequest req)(api.post="/book/add"),
    UpdateBookResponse updateBook(1: UpdateBookRequest req)(api.put="/book/update"),
    DeleteBookResponse deleteBook(1: DeleteBookRequest req)(api.delete="/book/delete"),
    GetBookResponse getBook(1: GetBookRequest req)(api.get="/book/get"),
}
