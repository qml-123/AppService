package main

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/qml-123/AppService/cgo/test"
	"github.com/qml-123/AppService/controller/file"
	"github.com/qml-123/AppService/controller/user"
	"github.com/qml-123/app_log/kitex_gen/base"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/kitex_gen/app"
)

// AppServiceImpl implements the last service interface defined in the IDL.
type AppServiceImpl struct{}

// Ping implements the AppServiceImpl interface.
func (s *AppServiceImpl) Ping(ctx context.Context, req *app.PingRequest) (resp *app.PingResponse, err error) {
	// TODO: Your code here...

	return &app.PingResponse{
		Message: "hello" + fmt.Sprintf("%d", test.Add(1, 2)),
	}, nil
}

// GetFile implements the AppServiceImpl interface.
func (s *AppServiceImpl) GetFile(ctx context.Context, req *app.GetFileRequest) (resp *app.GetFileResponse, err error) {
	// TODO: Your code here...
	return
}

// Upload implements the AppServiceImpl interface.
func (s *AppServiceImpl) Upload(ctx context.Context, req *app.UploadFileRequest) (resp *app.UploadFileResponse, err error) {
	// TODO: Your code here...
	file_key, err := file.Upload(ctx, req.GetUserId(), req.GetFile())
	if err != nil {
		if bizErr, ok := err.(*error_code.StatusError); ok {
			return &app.UploadFileResponse{
				BaseData: &base.BaseData{
					Code: thrift.Int32Ptr(int32(bizErr.Code)),
					Message: thrift.StringPtr(bizErr.Message),
				},
			}, nil
		}
		return nil, err
	}

	return &app.UploadFileResponse{
		FileKey: file_key,
	}, nil
}

// Login implements the AppServiceImpl interface.
func (s *AppServiceImpl) Login(ctx context.Context, req *app.LoginRequest) (resp *app.LoginResponse, err error) {
	// TODO: Your code here...
	if req.UserName == nil || req.Password == nil {
		return nil, error_code.InvalidParam
	}
	userID, err := user.Login(ctx, req)
	if err != nil {
		if bizErr, ok := err.(*error_code.StatusError); ok {
			return &app.LoginResponse{
				BaseData: &base.BaseData{
					Code: thrift.Int32Ptr(int32(bizErr.Code)),
					Message: thrift.StringPtr(bizErr.Message),
				},
			}, nil
		}
		return nil, err
	}
	return &app.LoginResponse{
		UserId: userID,
	}, nil
}

// Register implements the AppServiceImpl interface.
func (s *AppServiceImpl) Register(ctx context.Context, req *app.RegisteRequest) (resp *app.RegisteResponse, err error) {
	// TODO: Your code here...
	err = user.Register(ctx, req)
	if err != nil {
		if bizErr, ok := err.(*error_code.StatusError); ok {
			return &app.RegisteResponse{
				BaseData: &base.BaseData{
					Code: thrift.Int32Ptr(int32(bizErr.Code)),
					Message: thrift.StringPtr(bizErr.Message),
				},
			}, nil
		}
		return nil, err
	}
	return &app.RegisteResponse{}, nil
}

// GetFileKey implements the AppServiceImpl interface.
func (s *AppServiceImpl) GetFileKey(ctx context.Context, req *app.GetFileKeyRequest) (resp *app.GetFileKeyResponse, err error) {
	// TODO: Your code here...
	file_key, err := file.GetFileKey(ctx, req.GetUserId())
	if err != nil {
		if bizErr, ok := err.(*error_code.StatusError); ok {
			return &app.GetFileKeyResponse{
				BaseData: &base.BaseData{
					Code: thrift.Int32Ptr(int32(bizErr.Code)),
					Message: thrift.StringPtr(bizErr.Message),
				},
			}, nil
		}
		return nil, err
	}
	return &app.GetFileKeyResponse{
		FileKey: file_key,
	}, nil
}
