package file

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/qml-123/AppService/pkg/av1"
	"github.com/qml-123/AppService/pkg/db"
	_func "github.com/qml-123/AppService/pkg/func"
	"github.com/qml-123/app_log/logger"
)

func getChunkSize(fileSize int64) int64 {
	const MB = 1 << 20 // Define the size of a megabyte
	const MAX_CHUNK_SIZE = 100 * MB

	// Set chunkSize proportional to fileSize, with a minimum and maximum
	chunkSize := int64((float64(fileSize) * 99.0 / 10240.0))

	// Ensure chunkSize is at least 1MB
	if chunkSize < MB {
		chunkSize = MB
	}

	if chunkSize > MAX_CHUNK_SIZE {
		chunkSize = MAX_CHUNK_SIZE
	}

	return chunkSize
}

func splitFileChunk(ctx context.Context, d *av1.DirFile, filename string) error {
	logger.Info(ctx, "splitFileChunk begin, file_key: %v, file_name: %v", d.FileKey, filename)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileTotalSize := fileInfo.Size()
	size := getChunkSize(fileTotalSize)
	nowTotalSize := int64(0)
	buffer := make([]byte, size)
	for chunkNum := 1; ; chunkNum++ {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		nowTotalSize += int64(bytesread)
		tmpBuffer := buffer[:bytesread]
		has_more := !(nowTotalSize == fileTotalSize)
		err = _func.RetryFunc(func(ctx context.Context) error {
			var count int64
			f := &db.File{
				FileKey:      d.FileKey,
				FileType:     "video",
				Chunk:        tmpBuffer,
				ChunkNum:     chunkNum,
				ChunkSize:    int(size),
				Delete:       false,
				OwnUserID:    d.UserID,
				HasMore:      has_more,
				IsCompressed: true,
			}
			result := db.GetDB().Model(&db.File{}).Select("id").Where("file_key = ? and chunk_num = ? and `delete` = ? and is_compressed = ?", d.FileKey, chunkNum, false, true).Count(&count)
			if result.Error != nil {
				return result.Error
			}
			if count == 0 {
				result = db.GetDB().Create(f)
				if result.Error != nil {
					return result.Error
				}
				return nil
			}
			result = db.GetDB().Where("file_key = ? and chunk_num = ? and `delete` = ?", d.FileKey, chunkNum, false).Updates(f)
			if result.Error != nil {
				return result.Error
			}
			return nil
		}, _func.NewRetryEntity(), ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func splitFile(inputPath string, chunkSize int64) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	buffer := make([]byte, chunkSize)
	var i int64 = 1
	for {
		bytesRead, err := inputFile.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		chunkPath := fmt.Sprintf("%s_chunk_%d", inputPath, i)
		i++
		err = writeChunk(chunkPath, buffer[:bytesRead])
		if err != nil {
			return err
		}
	}
	return nil
}

func writeChunk(chunkPath string, buf []byte) error {
	chunkFile, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	_, err = chunkFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func contactChunk(ctx context.Context, file_key string, file_name string) error {
	file_struct := &db.File{}
	result := db.GetDB().Select("chunk, chunk_num, chunk_size, file_type, user_id").First(file_struct, "file_key = ? and `delete` = ? and is_compressed = ? and has_more = false", file_key, false, false)
	if result.Error != nil {
		return result.Error
	}
	var count int64
	result = db.GetDB().Model(file_struct).Select("id").Where("file_key = ? and `delete` = ? and is_compressed = ?", file_key, false, false).Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if int(count) != file_struct.ChunkNum {
		return fmt.Errorf("the chunks is not equal max chunk_num, len = %d, max = %d", count, file_struct.ChunkNum)
	}
	file, err := os.Create(file_name)
	if err != nil {
		return err
	}
	defer file.Close()
	for i := 1; i <= int(count); i++ {
		result = db.GetDB().Select("chunk").First(file_struct, "file_key = ? and `delete` = ? and is_compressed = ?", file_key, false, false)
		if result.Error != nil {
			return result.Error
		}
		_, err = file.Write(file_struct.Chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func ContactFileChunks(ctx context.Context, file_keys []string) {
	dirPath := "/opt/app/"
	for _, file_key := range file_keys {
		tmpDir := dirPath + file_key + "/"
		err := os.MkdirAll(tmpDir, os.ModePerm)
		if err != nil {
			logger.Warn(ctx, "file_key(%s) mkdir failed, err: %v", file_key, err)
			continue
		}
		err = contactChunk(ctx, file_key, tmpDir + file_key)
		if err != nil {
			logger.Warn(ctx, "file_key(%s) mkdir failed, err: %v", file_key, err)
			continue
		}
	}
}
