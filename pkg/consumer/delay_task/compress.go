package delay_task

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qml-123/AppService/pkg/file"
	"github.com/qml-123/app_log/logger"
)

func CompressHandleMessage(ctx context.Context, msg []byte) error {
	var file_keys []string
	if err := json.Unmarshal(msg, &file_keys); err != nil {
		return err
	}

	if len(file_keys) == 0 {
		return nil
	}

	go func() {
		failedFileKeys := file.Compress(ctx, file_keys, "/opt/app/")
		// failed 延时任务
		body, err := json.Marshal(failedFileKeys)
		if err != nil {
			logger.Warn(ctx, "Marshal failed, failedFileKeys: %v, err: %v", failedFileKeys, err)
			return
		}
		err = AddTask("compress", body, 10 * time.Second)
		if err != nil {
			logger.Warn(ctx, "AddTask failed, err: %v", err)
			return
		}
	}()
	return nil
}
