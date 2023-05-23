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
    1: required i64 user_id
    2: required binary file
    3: required i32 chunk_num
    4: required i32 chunk_size
    5: required bool has_more
    6: required string file_key
    7: required string file_type

    255: base.BaseData baseData
}

struct UploadFileResponse {
    255: base.BaseData baseData
}

struct GetFileChunkNumRequest {
    1: required i64 user_id
    2: required string file_key

    255: base.BaseData baseData
}

struct GetFileChunkNumResponse {
    1: required i32 chunk_num

    255: base.BaseData baseData
}

struct GetFileRequest {
    1: required i64 user_id
    2: required string file_key
    3: required i32 chunk_num
    4: optional i32 chunk_size

    255: base.BaseData baseData
}

struct GetFileResponse {
    1: required binary file
    2: required string file_type
    3: required bool has_more
    4: required i32 chunk_size
    5: required i32 total_num

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
  GetFileChunkNumResponse GetFileChunkSize(1: GetFileChunkNumRequest req)

  RegisteResponse Register(1: RegisteRequest req)
  LoginResponse Login(1: LoginRequest req)
}
