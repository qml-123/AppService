package file

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/qml-123/AppService/pkg/av1"
	"github.com/qml-123/AppService/pkg/db"
	_func "github.com/qml-123/AppService/pkg/func"
	"github.com/qml-123/app_log/logger"
	"gorm.io/gorm"
)

func Compress(ctx context.Context, file_keys []string, dirPath string) (failedFileKeys []string) {
	logger.Info(ctx, "file Compress begin, file_keys: %v", file_keys)

	failedFileKeys = make([]string, 0)
	files := make([]*av1.DirFile, 0)
	file_struct := &db.File{}
	for _, file_key := range file_keys {
		{
			var count int64
			result := db.GetDB().Model(&db.File{}).Select("id").Where("file_key = ? and `delete` = ? and is_compressed = ? and has_more = ?", file_key, false, true, false).Count(&count)
			if result.Error != nil {
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					failedFileKeys = append(failedFileKeys, file_key)
				}
				logger.Warn(ctx, "Compress db error, file_key: %v, err: %v", file_key, result.Error)
				continue
			}
			if count > 0 {
				continue
			}
		}
		{
			result := db.GetDB().Select("file_type").First(file_struct, "file_key = ? and `delete` = ?", file_key, false)
			if result.Error != nil {
				if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
					failedFileKeys = append(failedFileKeys, file_key)
				}
				logger.Warn(ctx, "Compress db error, file_key: %v, err: %v", file_key, result.Error)
				continue
			}
			// 非视频不处理
			if file_struct.FileType != db.VideoColumnValue {
				continue
			}
		}
		tmpDir := dirPath + file_key + "/"
		err := os.MkdirAll(tmpDir, os.ModePerm)
		if err != nil {
			logger.Warn(ctx, "mkdir(%s) err: %v", tmpDir, err)
			failedFileKeys = append(failedFileKeys, file_key)
			continue
		}

		result := db.GetDB().Select("chunk, chunk_num, chunk_size, file_type, user_id").First(file_struct, "file_key = ? and `delete` = ? and is_compressed = ? and has_more = false", file_key, false, false)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				failedFileKeys = append(failedFileKeys, file_key)
			}
			logger.Warn(ctx, "solve file(%s) failed", file_key)
			continue
		}
		var tsFiles []string
		{
			lines := strings.Split(string(file_struct.Chunk), "\n")
			for _, line := range lines {
				if strings.HasSuffix(line, ".ts") {
					tsFiles = append(tsFiles, line)
				}
			}
		}
		if len(tsFiles) != file_struct.ChunkNum-1 {
			logger.Warn(ctx, "tsFile len is not equal chunk nums, file_key: %v, len: %d, chunks: %d", file_key, len(tsFiles), file_struct.ChunkNum-1)
			continue
		}
		file_name := tmpDir + file_key + ".m3u8"
		file, err := os.Create(file_name)
		if err != nil {
			logger.Warn(ctx, "not create file(%s)", file_name)
			failedFileKeys = append(failedFileKeys, file_key)
			continue
		}
		_, err = file.Write(file_struct.Chunk)
		if err != nil {
			logger.Warn(ctx, "write file error, err: %v", err)
			continue
		}
		_ = file.Close()

		var userID int64
		for num := 1; num <= len(tsFiles); num++ {
			result = db.GetDB().Select("chunk, chunk_size, has_more, file_type, user_id").First(file_struct, "file_key = ? and `delete` = ? and is_compressed = ? and chunk_num = ?", file_key, false, false, num)
			if result.Error != nil {
				logger.Warn(ctx, "do sql failed, err: %v", err)
				break
			}
			file_name = tmpDir + tsFiles[num-1]
			file, err = os.Create(file_name)
			if err != nil {
				logger.Warn(ctx, "not create file(%s)", file_name)
				break
			}
			_, err = file.Write(file_struct.Chunk)
			if err != nil {
				logger.Warn(ctx, "write file error, err: %v", err)
				break
			}
			_ = file.Close()
		}
		if err != nil {
			failedFileKeys = append(failedFileKeys, file_key)
			logger.Warn(ctx, "solve file(%s) failed", file_key)
			continue
		}

		files = append(files, &av1.DirFile{
			Dir:     tmpDir,
			File:    tmpDir + file_key + ".m3u8",
			FileKey: file_key,
			UserID:  userID,
		})
	}
	logger.Info(ctx, "len(files) is %d", len(files))

	type funcEntity struct {
		fn func(ctx context.Context, d *av1.DirFile) error
		e  *_func.RetryEntity
		d  *av1.DirFile
	}

	type FuncTask struct {
		Task []*funcEntity
		Add  func(*funcEntity)
		Wait func(context.Context) error
	}
	t := &FuncTask{Task: make([]*funcEntity, 0)}
	t.Add = func(e *funcEntity) {
		t.Task = append(t.Task, e)
	}
	t.Wait = func(ctx context.Context) error {
		for _, task := range t.Task {
			if err := task.fn(ctx, task.d); err != nil {
				return err
			}
		}
		return nil
	}

	for i := 0; i < len(files); i++ {
		t.Add(&funcEntity{
			fn: av1.EncodeSplitFile,
			e:  _func.NewRetryEntity(),
			d:  files[i],
		})
		t.Add(&funcEntity{
			fn: av1.EncodeEachFile,
			e:  _func.NewRetryEntity(),
			d:  files[i],
		})
		t.Add(&funcEntity{
			fn: av1.Concatenate,
			e:  _func.NewRetryEntity(),
			d:  files[i],
		})
		if err := t.Wait(ctx); err != nil {
			failedFileKeys = append(failedFileKeys, files[i].FileKey)
			logger.Warn(ctx, "encode failed, err: %v", err)
			continue
		}
		filename := files[i].Dir + "final_output.m3u8"
		if err := splitFileChunk(ctx, files[i], filename); err != nil {
			failedFileKeys = append(failedFileKeys, files[i].FileKey)
			logger.Warn(ctx, "SplitFileChunk failed, err: %v", err)
			continue
		}
	}

	// clear files
	defer func() {
		for i := 0; i < len(files); i++ {
			if err := av1.DeleteGeneratedFiles(ctx, files[i]); err != nil {
				logger.Warn(ctx, "DeleteGeneratedFiles failed, err: %v", err)
			}
		}
	}()
	return
}
