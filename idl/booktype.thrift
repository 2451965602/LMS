namespace go booktype
include "model.thrift"

struct AddBookTypeRequest{
    1: required string title,
    2: required string author,
    3: required string category,
    4: required string ISBN,
    5: required string publisher,
    6: required string publishyear,
    7: required string description,
}
struct AddBookTypeResponse{
    1: model.BaseResp base,
    2: required i64 book_id,
}

struct UpdateBookTypeRequest{
    1: optional string title,
    2: optional string author,
    3: optional string category,
    4: required string ISBN,
    5: optional string publisher,
    6: optional string publishyear,
    7: optional string description,
}
struct UpdateBookTypeResponse{
    1: model.BaseResp base,
    2: required model.BookType data,
}

struct DeleteBookTypeRequest{
    1: required i64 ISBN,
}
struct DeleteBookTypeResponse{
    1: model.BaseResp base,
}

struct GetBookTypeRequest{
    1: optional i64 ISBN,
    2: optional string title,
    3: optional string author,
}
struct GetBookTypeResponse{
    1: model.BaseResp base,
    2: required list<model.BookType> data,
}


service BookTypeService {
    AddBookTypeResponse addBookType(1: AddBookTypeRequest req)(api.post="/booktype/add"),
    UpdateBookTypeResponse updateBookType(1: UpdateBookTypeRequest req)(api.put="/booktype/update"),
    DeleteBookTypeResponse deleteBookType(1: DeleteBookTypeRequest req)(api.delete="/booktype/delete"),
    GetBookTypeResponse getBookType(1: GetBookTypeRequest req)(api.get="/booktype/get"),
}