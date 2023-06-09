package file

import (
	"context"
	"errors"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/AppService/pkg/id"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

const (
	FileOwner  int = 1
	FileEditer int = 2
	FileReader int = 4
)

func isValidPermission(permission int) bool {
	return permission == FileOwner || permission == FileEditer || permission == FileReader
}

func GetFileKey(ctx context.Context, userID int64) (file_key string, err error) {
	file := &db.FileInfo{}
	for {
		file_key = id.GenerateFileKey()
		file.FileKey = file_key
		result := db.GetDB().Select("id").First(file, "file_key = ?", file_key)
		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			break
		}
		if result.Error != nil {
			logger.Error(ctx, "db First error, err: %v", result.Error)
			return "", err
		}
	}
	file.OwnUserID = userID
	file.FileKey = file_key
	file.Delete = false
	file.UploadEnd = false
	file.CompressionType = "av1"
	result := db.GetDB().Create(file)
	if result.Error != nil {
		logger.Error(ctx, "db Create error, err: %v", result.Error)
		return "", error_code.InternalError
	}
	return file_key, nil
}
