package user

import (
	"context"
	"errors"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/kitex_gen/app"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func Login(ctx context.Context, req *app.LoginRequest) (int64, error) {
	user := &db.User{}
	result := db.GetDB().Select("user_id").First(user, "user_name = ? and `delete` = ?", req.GetUserName(), false)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Info(ctx, "user is not exist, user_id: %v", req.GetUserName())
			return 0, error_code.NewStatus(error_code.InvalidParam.Code, "user is not exist")
		}
		logger.Warn(ctx, "db error, %v", result.Error)
		return 0, error_code.InternalError
	}
	return user.UserID, nil
}
