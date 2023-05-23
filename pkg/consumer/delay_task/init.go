package delay_task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/qml-123/AppService/pkg/file"
	"github.com/qml-123/AppService/pkg/redis"
	"github.com/qml-123/AppService/pkg/utils"
	"github.com/qml-123/app_log/logger"
)

var (
	handles = map[string]func(ctx context.Context, body []byte) error{
		"test":     TestHandleMessage,
		"compress": CompressHandleMessage,
	}
)

func Init() error {
	go func() {
		ctx := context.Background()
		for {
			result, err := redis.ZPopMinTask()
			if err != nil {
				logger.Info(ctx, "ZPopMinTask failed, err: %v", err)
				continue
			}
			if len(result) == 0 {
				time.Sleep(5 * time.Second)
				continue
			}
			str := result[0].Member.(string)
			strs := strings.Split(str, ":")
			if len(strs) != 2 {
				logger.Info(ctx, "jobID error, str : %s", str)
				continue
			}
			jobkey := strs[0]
			f, ok := handles[jobkey]
			if !ok {
				logger.Warn(ctx, "no handle found")
				continue
			}
			score := result[0].Score
			duration := time.Until(floatToTime(score))
			args, err := redis.HGetTask(str)
			if err != nil {
				logger.Warn(ctx, "HGetTask failed, str: %v", str)
				continue
			}
			if duration <= 0 {
				_ = f(ctx, []byte(args))
			} else {
				time.Sleep(duration)
				_ = f(ctx, []byte(args))
			}

			redis.HDelTask(str)
		}
	}()

	go func() {
		ctx := context.Background()
		for {
			time.Sleep(5 * time.Second)
			result, err := redis.SMember("file_keys")
			if err != nil {
				continue
			}
			var failedKeys []string
			file_keys, not_exists := check_valid_file_keys(ctx, result)
			//logger.Info(ctx, "not exists keys: %v", not_exists)
			if len(file_keys) > 0 {
				failedKeys = file.Compress(ctx, file_keys, "/opt/app/")
			}
			successKeys := make(map[string]bool)
			{
				m := utils.StrSliceToMap(failedKeys)
				for _, key := range file_keys {
					if _, ok := m[key]; ok {
						continue
					}
					successKeys[key] = true
				}
			}
			file_keys = append(utils.StrMapToSlice(successKeys), not_exists...)
			if len(file_keys) == 0 {
				continue
			}
			_, err = redis.SRem("file_keys", file_keys...)
			if err != nil {
				time.Sleep(2 * time.Second)
				if _, err = redis.SRem("file_keys", file_keys...); err != nil {
					logger.Warn(ctx, "redis SRem failed, file_keys: %v, err: %v", file_keys, err)
				}
			}
		}
	}()
	return nil
}

func AddTask(taskKey string, args []byte, delay time.Duration) error {
	jobID := fmt.Sprintf("%s:%d", taskKey, time.Now().UnixNano())
	if ok, err := redis.HSetTask(jobID, args); err != nil || !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("HSetTask ret false")
	}

	score := float64(time.Now().Add(delay).Unix())
	_, err := redis.ZAddTask(score, jobID)
	if err != nil {
		_, _ = redis.HDelTask(jobID)
		return err
	}
	return nil
}

func floatToTime(timestamp float64) time.Time {
	sec := int64(timestamp)
	nsec := int64((timestamp - float64(sec)) * 1e9)
	return time.Unix(sec, nsec)
}
