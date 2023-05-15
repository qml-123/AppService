package user

import (
	"context"
	"errors"

	"github.com/qml-123/AppService/kitex_gen/app"
	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func Register(ctx context.Context, req *app.RegisteRequest) error {
	var user *db.User
	result := db.GetDB().First(user, "user_name = ?", req.GetUserName())
	if result.Error == nil {
		logger.Warn(ctx, "the user_name is registered")
		return error_code.RegisterNameDuplicate
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Error(ctx, "db error, %v", result.Error)
		return error_code.InternalError
	}

	user = &db.User{
		UserName: req.GetUserName(),
		PassWord: req.GetPassword(),
	}
	result = db.GetDB().Create(user)
	if result.Error != nil {
		logger.Error(ctx, "db error, %v", result.Error)
		return error_code.InternalError
	}
	return nil
}
