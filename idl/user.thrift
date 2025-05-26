namespace go user
include "model.thrift"

struct RegisterRequest{
    1: required string username,
    2: required string password,
    3: required string phone,
}
struct RegisterResponse{
    1: model.BaseResp base,
    2: required i64 user_id,
}

struct LoginRequest{
    1: required string username,
    2: required string password,
}
struct LoginResponse{
    1: model.BaseResp base,
    2: required model.User data,
}

struct UpdateUserRequest{
    1: optional string phone,
    2: optional string password,
}
struct UpdateUserResponse{
    1: model.BaseResp base,
    2: required model.User data,
}

struct UserInfoRequest{
    1: required i64 user_id,
}
struct UserInfoResponse{
    1: model.BaseResp base,
    2: required model.User data,
}

struct DeleteUserRequest{
    1: required string username,
}
struct DeleteUserResponse{
    1: model.BaseResp base,
}

struct AdminUpdateUserRequest{
    1: required i64 user_id,
    2: optional string username,
    3: optional string phone,
    4: optional string permission,
    5: optional string status,
    6: optional string password,
}
struct AdminUpdateUserResponse{
    1: model.BaseResp base,
    2: required model.User data,
}

struct AdminDeleteUserRequest{
    1: required i64 user_id,
}
struct AdminDeleteUserResponse{
    1: model.BaseResp base,
}

struct RefreshTokenRequest{

}
struct RefreshTokenResponse{
    1: model.BaseResp base,
}


service UserService {
    RegisterResponse register(1: RegisterRequest req)(api.post="/user/register"),
    LoginResponse login(1: LoginRequest req)(api.post="/user/login"),
    UpdateUserResponse updateUser(1: UpdateUserRequest req)(api.put="/user/update"),
    UserInfoResponse GetUserInfo(1: UserInfoRequest req)(api.get="/user/info"),
    DeleteUserResponse deleteUser(1: DeleteUserRequest req)(api.delete="/user/delete"),
    RefreshTokenResponse refreshToken(1: RefreshTokenRequest req)(api.post="/user/refresh"),
}

service AdminUserService {
    AdminUpdateUserResponse adminUpdateUser(1: AdminUpdateUserRequest req)(api.put="/user/admin/update"),
    AdminDeleteUserResponse adminDeleteUser(1: AdminDeleteUserRequest req)(api.delete="/user/admin/delete"),
}
