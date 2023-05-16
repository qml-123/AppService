include "base.thrift"
namespace go app

struct PingRequest {
}

struct PingResponse {
  1: required string message

  255: base.BaseData baseData
}

struct GetFileKeyRequest {
    1: required i64 user_id

    255: base.BaseData baseData
}

struct GetFileKeyResponse {
    1: required string file_key

    255: base.BaseData baseData
}

struct UploadFileRequest {
    1: required string user_id
    2: required binary file
    255: base.BaseData baseData
}

struct UploadFileResponse {
    1: required string file_key

    255: base.BaseData baseData
}

struct GetFileRequest {
    1: required string user_id
    2: required string file_key

    255: base.BaseData baseData
}

struct GetFileResponse {
    1: required binary file

    255: base.BaseData baseData
}

struct RegisteRequest {
    1: required string user_name
    2: required string password
    3: optional string email
    4: optional string phone_number

    255: base.BaseData baseData
}

struct RegisteResponse {
    255: base.BaseData baseData
}

struct LoginRequest {
    1: optional string user_name
    2: optional string password
    3: optional string token

    255: base.BaseData baseData
}

struct LoginResponse {
    1: required string token
    2: required i64 user_id

    255: base.BaseData baseData
}

service AppService {
  PingResponse Ping(1: PingRequest req)

  GetFileResponse GetFile(1: GetFileRequest req)
  UploadFileResponse Upload(1: UploadFileRequest req)
  GetFileKeyResponse GetFileKey(1: GetFileKeyRequest req)

  RegisteResponse Register(1: RegisteRequest req)
  LoginResponse Login(1: LoginRequest req)
}
