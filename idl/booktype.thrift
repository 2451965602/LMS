namespace go booktype
include "model.thrift"

struct AddBookTypeRequest{
    1: required string title,
    2: required string author,
    3: required string category,
    4: required string ISBN,
    5: required string publisher,
    6: required i64 publish_year,
    7: required string description,
}
struct AddBookTypeResponse{
    1: model.BaseResp base,
    2: required model.BookType data,
}

struct UpdateBookTypeRequest{
    1: optional string title,
    2: optional string author,
    3: optional string category,
    4: required string ISBN,
    5: optional string publisher,
    6: optional i64 publish_year,
    7: optional string description,
}
struct UpdateBookTypeResponse{
    1: model.BaseResp base,
    2: required model.BookType data,
}

struct DeleteBookTypeRequest{
    1: required string ISBN,
}
struct DeleteBookTypeResponse{
    1: model.BaseResp base,
}

struct GetBookTypeRequest{
    1: optional string ISBN,
    2: optional string title,
    3: optional string author,
    4: optional string category,
    5: required i64 page_size,
    6: required i64 page_num,
}
struct GetBookTypeResponse{
    1: model.BaseResp base,
    2: required list<model.BookType> data,
    3: required i64 total,
}


service BookTypeService {
    AddBookTypeResponse addBookType(1: AddBookTypeRequest req)(api.post="/booktype/add"),
    UpdateBookTypeResponse updateBookType(1: UpdateBookTypeRequest req)(api.put="/booktype/update"),
    DeleteBookTypeResponse deleteBookType(1: DeleteBookTypeRequest req)(api.delete="/booktype/delete"),
    GetBookTypeResponse getBookType(1: GetBookTypeRequest req)(api.get="/booktype/get"),
}
