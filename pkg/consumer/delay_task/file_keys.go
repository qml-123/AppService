package delay_task

import (
	"context"
	"errors"

	"github.com/qml-123/AppService/pkg/db"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func check_valid_file_keys(ctx context.Context, keys []string) ([]string, []string) {
	var file_keys, not_exists []string

	file := &db.File{}
	for _, key := range keys {
		var total int64
		result := db.GetDB().Model(&db.File{}).Where("file_key = ? and `delete` = ? and is_compressed = ?", key, false, false).Count(&total)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				not_exists = append(not_exists, key)
			}
			continue
		}
		result = db.GetDB().Select("chunk_num").First(file, "file_key = ? and `delete` = ? and has_more = ? and is_compressed = ?", key, false, false, false)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				not_exists = append(not_exists, key)
			}
			continue
		}
		if int(total) != file.ChunkNum {
			logger.Info(ctx, "chunk total length is not equal chunkNum, count: %d, chunkNum: %d", total, file.ChunkNum)
			continue
		}
		file_keys = append(file_keys, key)
	}
	return file_keys, not_exists
}
