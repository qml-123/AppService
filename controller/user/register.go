package user

import (
	"context"
	"errors"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/AppService/pkg/id"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/kitex_gen/app"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func Register(ctx context.Context, req *app.RegisteRequest) error {
	user := &db.User{}
	result := db.GetDB().Select("id").First(user, "user_name = ? and `delete` = ?", req.GetUserName(), false)
	if result.Error == nil {
		logger.Warn(ctx, "the user_name is registered")
		return error_code.RegisterNameDuplicate
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Error(ctx, "db First error, err: %v", result.Error)
		return error_code.InternalError
	}
	userID, err := getUserID(ctx)
	if err != nil {
		return error_code.InternalError
	}
	user = &db.User{
		UserID:   userID,
		UserName: req.GetUserName(),
		PassWord: req.GetPassword(),
		Delete: false,
	}
	result = db.GetDB().Create(user)
	if result.Error != nil {
		logger.Error(ctx, "db Create error, err: %v", result.Error)
		return error_code.InternalError
	}
	return nil
}

func getUserID(ctx context.Context) (userID int64, err error) {
	userID = id.Generate().Int64()
	for {
		user := &db.User{}
		result := db.GetDB().Select("id").First(user, "user_id = ?", userID)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return
		}
		if result.Error != nil {
			logger.Error(ctx, "db First error, err: %v", result.Error)
			return 0, err
		}
	}
}
