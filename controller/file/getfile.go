package file

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func GetFile(ctx context.Context, userID int64, file_key string, chunk_num int32) ([]byte, string, int32, bool, int32, error) {
	file_info := &db.FileInfo{}
	result := db.GetDB().Select("user_id, upload_end").First(file_info, "file_key = ? and `delete` = ?", file_key, false)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Warn(ctx, "db First failed, err: %v", result.Error)
			return nil, "", 0, false, 0, error_code.NoPermission.WithErrMsg(" or file not exist")
		}
		logger.Warn(ctx, "db First failed, err: %v", result.Error)
		return nil, "", 0, false, 0, result.Error
	}
	if file_info.OwnUserID != userID {
		return nil, "", 0, false, 0, error_code.NoPermission
	}

	if !file_info.UploadEnd {
		return nil, "", 0, false, 0, error_code.FileNotEnd
	}

	file := &db.File{}
	result = db.GetDB().Select("chunk, chunk_num, file_type, chunk_size, has_more").First(file, "file_key = ? and chunk_num = ? and `delete` = ? and is_compressed = ?", file_key, chunk_num, false, false)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, "", 0, false, 0, error_code.NewStatus(error_code.InvalidParam.Code, "file not exist")
		}
		logger.Warn(ctx, "db First failed, err: %v", result.Error)
		return nil, "", 0, false, 0, result.Error
	}

	file_content := base64.StdEncoding.EncodeToString(file.Chunk)
	return []byte(file_content), file.FileType, int32(file.ChunkSize), file.HasMore, 0, nil
}

func GetFileChunkSize(ctx context.Context, user_id int64, file_key string) (int32, error) {
	file := &db.FileInfo{}
	result := db.GetDB().Select("user_id").First(file, "file_key = ? and `delete` = ?", file_key, false)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, error_code.NewStatus(error_code.InvalidParam.Code, "file not exist")
		}
		logger.Warn(ctx, "db first failed, file_key: %v, err: %v", file_key, result.Error)
		return 0, result.Error
	}

	var total int64
	result = db.GetDB().Model(&db.File{}).Where("file_key = ? and `delete` = ?", file_key, false).Count(&total)
	if result.Error != nil {
		logger.Warn(ctx, "db first failed, file_key: %v, err: %v", file_key, result.Error)
		return 0, result.Error
	}

	return int32(total), nil
}
