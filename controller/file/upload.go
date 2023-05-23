package file

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/app_log/error_code"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func Upload(ctx context.Context, user_id int64, file_content []byte, chunk_num, chunk_size int32, file_key string, has_more bool, file_type string) (err error) {
	file_info := &db.FileInfo{}
	result := db.GetDB().Select("user_id, upload_end").Where("file_key = ? and `delete` = ?", file_key, false).First(file_info)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return error_code.NewStatus(error_code.InvalidParam.Code, "file not exist")
		}
		return result.Error
	}

	if user_id != file_info.OwnUserID {
		return error_code.NoPermission
	}

	if file_info.UploadEnd {
		return error_code.FileExist
	}

	file_content, err = base64.StdEncoding.DecodeString(string(file_content))
	if err != nil {
		logger.Warn(ctx, "base64 Decode failed, err: %v", err)
		return err
	}

	file := &db.File{
		FileKey:      file_key,
		Chunk:        file_content,
		ChunkNum:     int(chunk_num),
		ChunkSize:    int(chunk_size),
		OwnUserID:    user_id,
		FileType:     file_type,
		HasMore:      has_more,
		Delete:       false,
		IsCompressed: false,
	}
	result = db.GetDB().Create(file)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate") {
			return error_code.FileExist
		}
		return result.Error
	}

	if !has_more {
		is_upload := false
		{
			var count int64
			result = db.GetDB().Model(file).Select("id").Where("file_key = ? and `delete` = false", file_key).Count(&count)
			if result.Error != nil {
				return result.Error
			}
			result = db.GetDB().Select("chunk_num").First(file, "file_key = ? and `delete` = false and has_more = false", file_key)
			if result.Error != nil {
				return result.Error
			}
			is_upload = int(count) == file.ChunkNum
		}
		if is_upload {
			result = db.GetDB().Model(&db.FileInfo{}).Where("file_key = ?", file_key).Update("upload_end", true)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	return nil
}
