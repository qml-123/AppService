package main

import (
	"context"
	"fmt"

	"github.com/qml-123/AppService/cgo/test"
	"github.com/qml-123/AppService/controller/file"
	"github.com/qml-123/AppService/controller/user"
	app "github.com/qml-123/AppService/kitex_gen/app"
	"github.com/qml-123/app_log/error_code"
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
		return nil, err
	}
	return &app.RegisteResponse{}, nil
}
